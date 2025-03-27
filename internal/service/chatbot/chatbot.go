package chatbot

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 新增配置文件对应结构体
type ChatBotRules struct {
	Metadata struct { // 新增元数据字段
		BotName string `mapstructure:"bot_name"`
		Version string `mapstructure:"version"`
	} `mapstructure:"metadata"`

	IntentDetection struct {
		RegexPatterns []struct {
			Intent        string   `mapstructure:"intent"`
			Patterns      []string `mapstructure:"patterns"`
			Priority      int      `mapstructure:"priority"`
			RequiredSlots []string `mapstructure:"required_slots"` // 新增字段
		} `mapstructure:"regex_patterns"`

		MLModel struct {
			Path      string  `mapstructure:"path"`
			Threshold float64 `mapstructure:"threshold"`
		} `mapstructure:"ml_model"` // 添加字段标签
	} `mapstructure:"intent_detection"` // 添加字段标签

	DialogueFlow struct {
		States []struct {
			Name         string       `mapstructure:"name"`
			Transitions  []Transition `mapstructure:"transitions"`
			EntryActions []Action     `mapstructure:"entry_actions"` // 新增字段
		} `mapstructure:"states"`
	} `mapstructure:"dialogue_flow"` // 添加字段标签

	ErrorHandling struct {
		DefaultFallback string `mapstructure:"default_fallback"`
	} `mapstructure:"error_handling"`
}

type State struct {
	Name        string       `mapstructure:"name"`
	Transitions []Transition `mapstructure:"transitions"`
}

type Transition struct {
	Intent    string   `mapstructure:"intent"`
	NextState string   `mapstructure:"next_state"` // 新增此字段
	Actions   []Action `mapstructure:"actions"`
}

type Action struct {
	Type    string                 `mapstructure:"type"`
	Content string                 `mapstructure:"content"`
	Key     string                 `mapstructure:"key"`
	Value   string                 `mapstructure:"value"`
	Params  map[string]interface{} `mapstructure:"params"`
}

// 新增查找状态的辅助方法
func (r *ChatBotRules) findState(name string) *State {
	for _, state := range r.DialogueFlow.States {
		if state.Name == name {
			return &State{
				Name:        state.Name,
				Transitions: state.Transitions,
			}
		}
	}
	return nil
}

type ChatBotEngine struct {
	db         *gorm.DB
	rules      ChatBotRules
	contextMap map[string]ConversationContext // key: customerID
}

type ConversationContext struct {
	CurrentState string
	Slots        map[string]string
	LastActive   time.Time
}

// 初始化聊天机器人
// 在ChatBotEngine结构体下补充缺失的方法
func (e *ChatBotEngine) getOrCreateContext(customerID string) ConversationContext {
	if ctx, exists := e.contextMap[customerID]; exists {
		return ctx
	}
	return ConversationContext{
		CurrentState: "welcome",
		Slots:        make(map[string]string),
		LastActive:   time.Now(),
	}
}

func (e *ChatBotEngine) saveContext(customerID string, ctx ConversationContext) {
	e.contextMap[customerID] = ctx
}

// 补充动作执行逻辑
func (e *ChatBotEngine) executeActions(actions []Action, ctx *ConversationContext) string {
	var response string
	for _, action := range actions {
		switch action.Type {
		case "response":
			response = action.Content // 简单实现，实际需要模板渲染
		case "set_context":

			if action.Key != "" {
				ctx.Slots[action.Key] = fmt.Sprintf("%v", action.Params["value"])
			} else {
				log.Printf("debug executeActions: %s, %v", "a", action)
			}
		}
	}
	return response
}

// 修改初始化方法加载配置
func NewChatBotEngine(db *gorm.DB) *ChatBotEngine {
	return &ChatBotEngine{
		db:         db,
		rules:      LoadChatBotRules(), // 加载配置
		contextMap: make(map[string]ConversationContext),
	}
}

// 在LoadChatBotRules函数中添加调试日志
func LoadChatBotRules() ChatBotRules {
	viper.SetConfigName("chatbot_rules")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	// 打印原始配置数据

	var rules ChatBotRules
	if err := viper.Unmarshal(&rules); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	return rules
}

func (e *ChatBotEngine) GetContext(customerID string) ConversationContext {
	return e.getOrCreateContext(customerID)
}

// 核心消息处理逻辑
func (e *ChatBotEngine) ProcessMessage(customerID string, message string) string {
	ctx := e.getOrCreateContext(customerID)
	defer e.saveContext(customerID, ctx)

	// 1. 意图识别
	intent := e.detectIntent(message, ctx)
	// 2. 状态转移
	response := e.handleStateTransition(intent, &ctx)

	// 3. 上下文更新
	ctx.LastActive = time.Now()
	return response
}

// 意图识别实现
func (e *ChatBotEngine) detectIntent(msg string, ctx ConversationContext) string {
	msg = strings.ToLower(msg)
	// 优先匹配正则规则
	for _, rule := range e.rules.IntentDetection.RegexPatterns {
		for _, pattern := range rule.Patterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(msg) {
				return rule.Intent
			}
		}
	}

	return "unknown"
}

// 状态机处理
// 修复空指针问题和状态转移逻辑
func (e *ChatBotEngine) handleStateTransition(intent string, ctx *ConversationContext) string {
	// 确保获取当前状态
	currentState := e.rules.findState(ctx.CurrentState)
	if currentState == nil {
		// 默认回退到欢迎状态
		currentState = e.rules.findState("welcome")
		if currentState == nil {
			return e.rules.ErrorHandling.DefaultFallback
		}
		ctx.CurrentState = currentState.Name // 更新上下文状态
	}

	// 增加空指针保护
	if currentState.Transitions == nil {
		return e.rules.ErrorHandling.DefaultFallback
	}

	// 优化状态转移匹配逻辑
	for _, transition := range currentState.Transitions {
		if transition.Intent == intent {
			response := e.executeActions(transition.Actions, ctx)
			if nextState := e.rules.findState(transition.NextState); nextState != nil {
				ctx.CurrentState = nextState.Name // 更新到下一个状态
			}
			return response
		}
	}
	return e.rules.ErrorHandling.DefaultFallback
}

/*
// 在现有聊天处理中集成
func ServeWebSocket(c *gin.Context, upgrader websocket.Upgrader) {
	db := c.MustGet("DB").(*gorm.DB)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	validCustomerID := uint64(1)
	chatbotEngine := NewChatBotEngine(db)
	for {
		// 接收消息
		_, msg, _ := conn.ReadMessage()

		// 处理业务逻辑
		response := chatbotEngine.ProcessMessage(
			strconv.FormatUint(validCustomerID, 10),
			string(msg),
		)

		// 发送响应
		conn.WriteMessage(websocket.TextMessage, []byte(response))
	}
}
*/
