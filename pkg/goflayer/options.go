package goflayer

// Options содержит конфигурацию для создания бота.
//
// Поля, которые имеют значения по умолчанию, указаны в комментариях.
// Используйте DefaultOptions() для получения преднастроенных опций.
//
// Пример использования:
//
//	opts := goflayer.Options{
//	    Host: "localhost",
//	    Port: 25565,
//	    Username: "MyBot",
//	    Version: "1.20.1",
//	}
//	bot, err := goflayer.CreateBot(opts)
type Options struct {
	// Обязательные параметры

	// Host - адрес сервера для подключения.
	// Может быть IP-адресом или доменным именем.
	Host string

	// Port - порт сервера Minecraft.
	// Обычно 25565 для серверов по умолчанию.
	Port int

	// Username - имя пользователя для бота.
	// В режиме offline будет использоваться как есть.
	// В онлайн режиме должен быть валидным именем учетной записи Microsoft.
	Username string

	// Параметры подключения

	// Version - версия Minecraft для подключения.
	// Если пустая, версия будет автоматически определена при подключении.
	// Поддерживаемые версии: см. константы LatestSupportedVersion, OldestSupportedVersion.
	Version string

	// Password - пароль для учетной записи Microsoft.
	// Обязателен только для онлайн режима с аутентификацией.
	// Используется вместе с Auth для входа в учетную запись.
	Password string

	// Auth - тип аутентификации.
	// "offline" - без аутентификации (для cracked/offline серверов)
	// "microsoft" - аутентификация через Microsoft Account (требует Password или код устройства)
	// "yggdrasil" - устаревшая аутентификация Mojang
	Auth string

	// SessionServer - кастомный сервер сессий (для advanced use cases).
	// Оставьте пустым для использования стандартных серверов.
	SessionServer string

	// AuthServer - кастомный сервер аутентификации.
	// Оставьте пустым для использования стандартных серверов.
	AuthServer string

	// Параметры подключения к прокси

	// Connect - функция для создания TCP соединения.
	// Может быть использована для подключения через прокси.
	// Если nil, используется стандартный net.Dial.
	//
	// Пример для SOCKS5 прокси:
	//	proxyDialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, &net.Dialer{})
	//	opts := goflayer.Options{
	//	    Host: "mc.server.com",
	//	    Port: 25565,
	//	    Username: "Bot",
	//	    Connect: func(ctx context.Context, network, address string) (net.Conn, error) {
	//	        return proxyDialer.Dial(network, address)
	//	    },
	//	}
	// Connect func(ctx context.Context, network, address string) (net.Conn, error)

	// Параметры клиента

	// Brand - идентификатор клиента, отправляемый серверу.
	// Обычно "vanilla" для стандартных клиентов.
	Brand string

	// Locale - язык клиента.
	// Используется для локализации сообщений от сервера.
	// По умолчанию "en_US".
	Locale string

	// Параметры поведения

	// Respawn - автоматически возрождаться после смерти.
	// Если true, бот будет автоматически респауниться.
	// По умолчанию true.
	Respawn bool

	// Параметры логирования

	// LogErrors - логировать ошибки в консоль.
	// Если true, ошибки будут выводиться через log.Print().
	// По умолчанию true.
	LogErrors bool

	// HideErrors - скрывать ошибки от вывода.
	// Если true, ошибки не будут выводиться в консоль.
	// По умолчанию false.
	HideErrors bool

	// Параметры плагинов

	// LoadInternalPlugins - загружать встроенные плагины.
	// Если true, все стандартные плагины будут загружены автоматически.
	// По умолчанию true.
	LoadInternalPlugins bool

	// Plugins - карта плагинов для загрузки или отключения.
	// Ключ - имя плагина, значение - false для отключения или функция плагина.
	//
	// Пример отключения плагина:
	//	opts := goflayer.Options{
	//	    // ...
	//	    Plugins: map[string]interface{}{
	//	        "physics": false,  // Отключить физический движок
	//	    },
	//	}
	//
	// Пример добавления внешнего плагина:
	//	opts := goflayer.Options{
	//	    // ...
	//	    Plugins: map[string]interface{}{
	//	        "myPlugin": &MyCustomPlugin{},
	//	    },
	//	}
	Plugins map[string]interface{}

	// Advanced параметры

	// Client - существующий протокол клиент для использования.
	// Если не nil, будет использован вместо создания нового.
	// Для advanced use cases когда нужно переиспользовать соединение.
	// Client interface{}

	// SkipValidation - пропустить валидацию версии сервера.
	// Если true, бот будет пытаться подключиться даже к неподдерживаемым версиям.
	// Опасно! Может привести к неожиданным ошибкам.
	// По умолчанию false.
	SkipValidation bool
}

// DefaultOptions возвращает опции с значениями по умолчанию.
// Используйте это как отправную точку и измените нужные поля.
//
// Пример:
//
//	opts := goflayer.DefaultOptions()
//	opts.Host = "mc.server.com"
//	opts.Username = "MyBot"
//	bot, err := goflayer.CreateBot(opts)
func DefaultOptions() Options {
	return Options{
		Host:                 "localhost",
		Port:                 25565,
		Username:             "Player",
		Version:              "", // Автоопределение
		Auth:                 "offline",
		Brand:                "vanilla",
		Locale:               "en_US",
		Respawn:              true,
		LogErrors:            true,
		HideErrors:           false,
		LoadInternalPlugins:  true,
		Plugins:              make(map[string]interface{}),
		SkipValidation:       false,
	}
}

// Validate проверяет опции на валидность и возвращает ошибку если невалидны.
func (o *Options) Validate() error {
	// Проверка обязательных полей
	if o.Host == "" {
		return ErrInvalidCredentials
	}
	if o.Port <= 0 || o.Port > 65535 {
		return ErrInvalidCredentials
	}
	if o.Username == "" {
		return ErrInvalidCredentials
	}

	// Проверка версии (если указана)
	// Здесь будет проверка на поддерживаемые версии
	// при добавлении системы версий

	// Проверка типа аутентификации
	if o.Auth == "" {
		o.Auth = "offline"
	}

	return nil
}

// Поддерживаемые версии Minecraft
const (
	// LatestSupportedVersion - последняя поддерживаемая версия Minecraft.
	LatestSupportedVersion = "1.21.5"

	// OldestSupportedVersion - самая старая поддерживаемая версия Minecraft.
	OldestSupportedVersion = "1.17"
)

// SupportFeature проверяет, поддерживается ли определенная функция в указанной версии.
// Это вспомогательная функция для проверки возможностей версий.
//
// Пример:
//
//	if goflayer.SupportFeature("blockSkull", "1.20") {
//	    // Черепа поддерживаются как блоки в 1.20+
//	}
func SupportFeature(feature, version string) bool {
	// TODO: Реализовать проверку фич на основе minecraft-data
	// Это будет добавлено при интеграции registry
	return false
}
