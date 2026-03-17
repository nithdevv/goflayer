// Package goflayer предоставляет основной API для создания Minecraft ботов.
//
# Базовые типы ошибок

package goflayer

import "errors"

// Базовые ошибки, которые могут возникать при работе с goflayer.
// Используйте errors.Is() и errors.As() для проверки ошибок.

var (
	// ErrBotNotConnected возникает при попытке выполнить операцию без соединения.
	ErrBotNotConnected = errors.New("bot: not connected to server")

	// ErrBotAlreadyConnected возникает при попытке подключиться, когда бот уже подключен.
	ErrBotAlreadyConnected = errors.New("bot: already connected to server")

	// ErrBotStopped возникает при попытке использовать остановленного бота.
	ErrBotStopped = errors.New("bot: bot is stopped")

	// ErrInvalidCredentials возникает при неправильных учетных данных для аутентификации.
	ErrInvalidCredentials = errors.New("auth: invalid credentials")

	// ErrServerOffline возникает при попытке подключиться к серверу в режиме онлайн с offline сервером.
	ErrServerOffline = errors.New("server: server is in offline mode")

	// ErrConnectionRefused возникает при отказе в соединении.
	ErrConnectionRefused = errors.New("connection: refused by server")

	// ErrConnectionTimeout возникает при таймауте соединения.
	ErrConnectionTimeout = errors.New("connection: timeout")

	// ErrDisconnected возникает при неожиданном разрыве соединения.
	ErrDisconnected = errors.New("connection: disconnected")

	// ErrPacketRead возникает при ошибке чтения пакета.
	ErrPacketRead = errors.New("protocol: packet read error")

	// ErrPacketWrite возникает при ошибке записи пакета.
	ErrPacketWrite = errors.New("protocol: packet write error")

	// ErrInvalidPacket возникает при получении невалидного пакета.
	ErrInvalidPacket = errors.New("protocol: invalid packet")

	// ErrUnknownPacket возникает при получении неизвестного пакета.
	ErrUnknownPacket = errors.New("protocol: unknown packet")

	// ErrCompressionError возникает при ошибке сжатия/расжатия.
	ErrCompressionError = errors.New("protocol: compression error")

	// ErrEncryptionError возникает при ошибке шифрования/дешифрования.
	ErrEncryptionError = errors.New("protocol: encryption error")

	// ErrInvalidState возникает при операции в неправильном состоянии протокола.
	ErrInvalidState = errors.New("protocol: invalid state for this operation")

	// ErrTimeout возникает при таймауте операции.
	ErrTimeout = errors.New("operation: timeout")

	// ErrCancelled возникает при отмене операции через контекст.
	ErrCancelled = errors.New("operation: cancelled")

	// ErrEntityNotFound возникает при попытке найти несуществующую сущность.
	ErrEntityNotFound = errors.New("entity: not found")

	// ErrBlockNotFound возникает при попытке найти несуществующий блок.
	ErrBlockNotFound = errors.New("block: not found")

	// ErrPlayerNotFound возникает при попытке найти несуществующего игрока.
	ErrPlayerNotFound = errors.New("player: not found")

	// ErrInventoryNotFound возникает при попытке открыть несуществующий инвентарь.
	ErrInventoryNotFound = errors.New("inventory: not found")

	// ErrItemNotFound возникает при попытке найти несуществующий предмет.
	ErrItemNotFound = errors.New("item: not found")

	// ErrSlotOutOfRange возникает при обращении к несуществующему слоту.
	ErrSlotOutOfRange = errors.New("inventory: slot out of range")

	// ErrWindowClosed возникает при попытке взаимодействовать с закрытым окном.
	ErrWindowClosed = errors.New("window: window is closed")

	// ErrTransactionRejected возникает при отклонении транзакции сервером.
	ErrTransactionRejected = errors.New("inventory: transaction rejected by server")

	// ErrDiggingCancelled возникает при отмене рытья блока.
	ErrDiggingCancelled = errors.New("digging: digging cancelled")

	// ErrCannotDig возникает когда блок нельзя копать.
	ErrCannotDig = errors.New("digging: cannot dig this block")

	// ErrTooFar возникает когда бот слишком далеко от цели.
	ErrTooFar = errors.New("action: too far from target")

	// ErrNotLookingAt возникает когда бот не смотрит на цель.
	ErrNotLookingAt = errors.New("action: not looking at target")

	// ErrPluginAlreadyLoaded возникает при попытке загрузить уже загруженный плагин.
	ErrPluginAlreadyLoaded = errors.New("plugin: already loaded")

	// ErrPluginLoadFailed возникает при ошибке загрузки плагина.
	ErrPluginLoadFailed = errors.New("plugin: load failed")

	// ErrVersionNotSupported возникает при попытке подключиться к неподдерживаемой версии.
	ErrVersionNotSupported = errors.New("version: not supported")

	// ErrRegistryNotLoaded возникает при попытке использовать незагруженный реестр.
	ErrRegistryNotLoaded = errors.New("registry: not loaded")

	// ErrChunkNotLoaded возникает при попытке обратиться к незагруженному чанку.
	ErrChunkNotLoaded = errors.New("world: chunk not loaded")

	// ErrPhysicsError возникает при ошибке физического движка.
	ErrPhysicsError = errors.New("physics: simulation error")

	// ErrPathNotFound возникает при отсутствии пути к цели.
	ErrPathNotFound = errors.New("pathfinding: no path found")

	// ErrInvalidPosition возникает при невалидной позиции.
	ErrInvalidPosition = errors.New("position: invalid position")

	// ErrInvalidDimension возникает при ошибке смены измерения.
	ErrInvalidDimension = errors.New("world: invalid dimension")

	// ErrNotInGame возникает при попытке выполнить действие вне игрового состояния.
	ErrNotInGame = errors.New("bot: not in game")
)
