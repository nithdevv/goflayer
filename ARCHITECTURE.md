# goflayer Architecture

This document provides a comprehensive overview of goflayer's architecture for developers and AI assistants continuing this project.

## Overview

goflayer is a complete rewrite of mineflayer (JavaScript) to Golang. This document captures all architectural decisions, implementation progress, and guidance for future development.

## Project Goals

1. **Feature Parity** - All 47 mineflayer plugins
2. **Performance** - Leverage Go's concurrency model
3. **Type Safety** - Strong typing over JavaScript's dynamic types
4. **Documentation** - English comments for universal understanding
5. **Maintainability** - Clean architecture for AI collaboration

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         User Code                           │
│  bot.Chat("Hello")  │  bot.On("chat", handler)              │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│                      Bot Interface                          │
│  - CreateBot()    - Connect()    - On()                     │
│  - Chat()         - MoveTo()     - LoadPlugin()             │
└────────────────────┬────────────────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         │                       │
┌────────▼────────┐    ┌────────▼─────────────┐
│  Plugin System  │    │  Protocol Client     │
│  - Loader       │    │  - Serialization     │
│  - Interfaces   │    │  - Compression       │
│  - 47 Plugins   │    │  - Encryption        │
└─────────────────┘    └──────────┬───────────┘
                                  │
                         ┌────────▼────────────┐
                         │  Network Layer       │
                         │  - TCP Connection    │
                         │  - Packet I/O        │
                         └─────────────────────┘
```

## Package Structure

### `pkg/goflayer/`
**Purpose**: Main bot interface and core types

**Key Files**:
- `bot.go` - Bot interface and implementation
- `options.go` - Configuration options
- `events.go` - Event system
- `errors.go` - Error definitions

**Important Types**:
```go
type Bot interface {
    Connect(ctx) error
    Disconnect() error
    On(event, handler) Subscription
    Chat(message) error
    // ... more methods
}
```

### `pkg/protocol/`
**Purpose**: Minecraft protocol implementation

**Key Files**:
- `client.go` - Protocol client
- `packet.go` - Packet structures
- `serializer.go` - Binary serialization
- `states.go` - Protocol state machine
- `compression.go` - Zlib compression
- `encryption.go` - AES encryption

**Protocol States**:
```go
const (
    Handshaking State = iota  // Initial connection
    Status                    // Server ping
    Login                     // Authentication
    Play                      // Gameplay
)
```

### `pkg/plugins/`
**Purpose**: Plugin system and implementations

**Key Files**:
- `interface.go` - Plugin interfaces
- `loader.go` - Plugin loader

**Plugin Interface**:
```go
type Plugin interface {
    Name() string
    Version() string
    Inject(bot, options) error
    Cleanup() error
}
```

### `pkg/math/`
**Purpose**: 3D vector mathematics

**Key File**: `vec3.go` - Complete Vec3 implementation

### `pkg/registry/`
**Purpose**: Version-specific data (TODO)

## Implementation Progress

### ✅ Completed (Phase 1)

#### Core Infrastructure
- [x] Project structure and go.mod
- [x] Vec3 math library (`pkg/math/vec3.go`)
- [x] Error types (`pkg/goflayer/errors.go`)
- [x] Bot options (`pkg/goflayer/options.go`)
- [x] Event system (`pkg/goflayer/events.go`)
- [x] Bot interface (`pkg/goflayer/bot.go`)

#### Protocol Layer
- [x] Protocol states (`pkg/protocol/states.go`)
- [x] Packet structure (`pkg/protocol/packet.go`)
- [x] Protocol client (`pkg/protocol/client.go`)
- [x] Serialization (`pkg/protocol/serializer.go`)
- [x] Zlib compression (`pkg/protocol/compression.go`)
- [x] AES encryption (`pkg/protocol/encryption.go`)

#### Plugin System
- [x] Plugin interfaces (`pkg/plugins/interface.go`)
- [x] Plugin loader (`pkg/plugins/loader.go`)

#### Documentation
- [x] README.md
- [x] ARCHITECTURE.md (this file)
- [x] Example bot (`cmd/example/main.go`)

### 🚧 In Progress (Phase 2)

#### Protocol Implementation
- [ ] Complete packet definitions
- [ ] Protocol version detection
- [ ] Handshake flow
- [ ] Login flow
- [ ] Play state packets

### 📋 Planned (Phase 3+)

#### Core Plugins (47 total)
1. `game.go` - Game state
2. `entities.go` - Entity tracking (935 lines in JS)
3. `physics.go` - Physics engine
4. `blocks.go` - Block management
5. `chat.go` - Chat system
6. `inventory.go` - Inventory handling
7. `digging.go` - Block digging
8. `place_block.go` - Block placement
9. `movement.go` - Movement controls
10. `craft.go` - Crafting system
11. `health.go` - Health tracking
12. `experience.go` - Experience system
13. And 34 more plugins...

See mineflayer's `lib/plugins/` for complete list.

## Key Architectural Patterns

### 1. Promise → Channel Conversion

**JavaScript (mineflayer)**:
```javascript
async function digBlock(block) {
  await bot.lookAt(block.position);
  await bot.dig(block);
  await sleep(1000);
}
```

**Go (goflayer)**:
```go
func (p *DiggingPlugin) DigBlock(ctx context.Context, block *Block) error {
    task := sync.NewTask[bool]()

    sub := p.bot.On("diggingCompleted", func(data ...interface{}) {
        task.Finish(true)
    })
    defer sub.Unsubscribe()

    if err := p.bot.LookAt(ctx, block.Position); err != nil {
        return err
    }

    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    _, err := task.Wait(ctx)
    return err
}
```

### 2. EventEmitter → EventBus

**JavaScript**:
```javascript
bot.on('chat', (username, message) => {
  console.log(`${username}: ${message}`);
});
```

**Go**:
```go
bot.On("chat", func(data ...interface{}) {
    username := data[0].(string)
    message := data[1].(string)
    fmt.Printf("%s: %s\n", username, message)
})
```

### 3. Plugin Injection

**JavaScript**:
```javascript
function plugin(bot, options) {
  bot.myMethod = function() { ... };
}
```

**Go**:
```go
type MyPlugin struct {
    bot goflayer.Bot
}

func (p *MyPlugin) Inject(bot goflayer.Bot, options goflayer.Options) error {
    p.bot = bot
    return nil
}
```

## File Reference Guide

### Critical Files from mineflayer

| Component | JS File | Go Equivalent | Status |
|-----------|---------|---------------|--------|
| Bot Factory | `lib/loader.js` | `pkg/goflayer/bot.go` | ✅ |
| Plugin Loader | `lib/plugin_loader.js` | `pkg/plugins/loader.go` | ✅ |
| Entities | `lib/plugins/entities.js` | `pkg/plugins/entities.go` | 📋 |
| Physics | `lib/plugins/physics.js` | `pkg/plugins/physics.go` | 📋 |
| Inventory | `lib/plugins/inventory.js` | `pkg/plugins/inventory.go` | 📋 |
| Chat | `lib/plugins/chat.js` | `pkg/plugins/chat.go` | 📋 |
| Protocol | `temp/mcp-protocol/src/client.js` | `pkg/protocol/client.go` | ✅ |
| Serialization | `temp/mcp-protocol/src/transforms/serializer.js` | `pkg/protocol/serializer.go` | ✅ |

## Implementation Guidelines

### When Converting JS → Go

1. **Async Functions**: Convert to functions accepting `context.Context`
2. **Promises**: Use channels or custom task types
3. **Callbacks**: Use Go channels or event subscriptions
4. **Dynamic Types**: Use interface{} with type assertions or define structs
5. **Event Handlers**: Use goroutine-safe event bus

### Go Best Practices

1. **Error Handling**: Always return errors, never panic in library code
2. **Concurrency**: Use goroutines for parallel processing, channels for communication
3. **Context**: Always accept and propagate context.Context
4. **Mutex**: Use sync.RWMutex for read-heavy data structures
5. **Documentation**: Export all types with godoc comments

### Performance Optimization

1. **Object Pooling**: Use `sync.Pool` for frequently allocated types (Vec3, Packet)
2. **Zero-Copy**: Use sub-slices instead of copying buffers
3. **Preallocation**: Preallocate maps/slices when size is known
4. **Caching**: Cache registry lookups, version checks

## Testing Strategy

### Unit Tests
```go
func TestVec3_Add(t *testing.T) {
    v1 := NewVec3(1, 2, 3)
    v2 := NewVec3(4, 5, 6)
    result := v1.Plus(v2)

    if result.X != 5 || result.Y != 7 || result.Z != 9 {
        t.Errorf("Expected (5, 7, 9), got %v", result)
    }
}
```

### Integration Tests
```go
func TestBot_Connect(t *testing.T) {
    // Start test server
    server := NewTestServer()
    defer server.Close()

    // Connect bot
    bot, err := CreateBot(Options{
        Host: server.Host,
        Port: server.Port,
    })
    assert.NoError(t, err)
    assert.NoError(t, bot.Connect(context.Background()))
}
```

## Version Support

**Currently Targeting**: Minecraft 1.17 - 1.21.5

**Versions**:
- Latest: 1.21.5
- Oldest: 1.17
- Primary: 1.20.1

## Dependencies

### External
- None (planned to be standalone)

### Internal
- `github.com/go-flayer/goflayer/pkg/protocol`
- `github.com/go-flayer/goflayer/pkg/math`
- `github.com/go-flayer/goflayer/pkg/plugins`

## Development Workflow

1. **Pick a plugin** from the TODO list
2. **Read the JS source** in `temp/mineflayer/lib/plugins/`
3. **Create Go file** in `pkg/plugins/`
4. **Implement interface** - Name(), Version(), Inject(), Cleanup()
5. **Convert logic** - Apply JS→Go patterns
6. **Add tests** - Unit and integration
7. **Update docs** - README, ARCHITECTURE.md

## Common Tasks

### Adding a New Plugin

```bash
# 1. Create file
touch pkg/plugins/myplugin.go

# 2. Implement interface
type MyPlugin struct { ... }

# 3. Load in bot.go
func (b *bot) loadInternalPlugins() error {
    b.LoadPlugin(&MyPlugin{})
    // ...
}
```

### Adding a Packet Handler

```go
// In plugin Inject()
bot.ProtocolClient().On("packet_name", func(packet *protocol.Packet) {
    data := packet.Data
    // Handle packet
})
```

### Emitting Events

```go
bot.Emit("myEvent", data1, data2)
```

## Troubleshooting

### Common Issues

1. **Import Errors**: Run `go mod tidy`
2. **Type Assertions**: Always check ok value
3. **Goroutine Leaks**: Always use context cancellation
4. **Race Conditions**: Use sync.RWMutex for shared state

### Debug Mode

Set environment variable:
```bash
DEBUG=1 go run cmd/example/main.go
```

## Performance Benchmarks

TODO: Add benchmarks comparing goflayer vs mineflayer

## Contributing

See CONTRIBUTING.md for detailed guidelines.

## License

MIT License - See LICENSE file

---

**Last Updated**: 2025-01-17
**Version**: 0.1.0 (Alpha)
**Status**: Active Development
