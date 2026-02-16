### Утилиты для получения данных о зарегистрированных телефонах на Asterisk с PJSIP по протоколу AMI

### Публичные функции

### GetPJSIPEndpointsIPtoDataMap

    func GetPJSIPEndpointsIPtoDataMap(AMIAddr string, AMIPort int, AMIUsername string, AMIPassword string) (map[string]ExtensionData, error)

Получает данные о всех зарегистрированных PJSIP endpoint'ах из Asterisk через AMI протокол.

**Параметры:**
- AMIAddr - IP адрес или hostname сервера Asterisk
- AMIPort - порт AMI (обычно 5038)
- AMIUsername - имя пользователя AMI
- AMIPassword - пароль AMI

**Возвращает:** карту данных о endpoint'ах, индексированную по IP адресу, или ошибку при неудаче.

**Процесс работы:**
1. Подключается к Asterisk AMI
2. Аутентифицируется
3. Запрашивает список всех PJSIP endpoint'ов
4. Для каждого endpoint'а получает детальную информацию (Caller ID)
5. Парсит SIP Contact URI для извлечения IP адресов
6. Возвращает данные, индексированные по IP

**Пример:**

    endpoints, err := asteriskamitools.GetPJSIPEndpointsIPtoDataMap(
        "10.0.0.1", 
        5038, 
        "admin", 
        "password",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Поиск телефона по IP
    if phone, exists := endpoints["10.0.0.50"]; exists {
        fmt.Printf("Телефон на IP 10.0.0.50:\n")
        fmt.Printf("  Extension: %s\n", phone.Extension)
        fmt.Printf("  Name: %s\n", phone.Name)
    }
    
    // Вывод всех зарегистрированных телефонов
    fmt.Printf("Всего зарегистрировано: %d телефонов\n", len(endpoints))

---

### Типы данных

#### ExtensionData

    type ExtensionData struct {
        Extension string // Номер внутреннего телефона (например, "100")
        Contacts  string // Полный SIP Contact URI (например, "sip:100@10.0.0.50:5060")
        IP        string // IP адрес, извлеченный из Contact (например, "10.0.0.50")
        Name      string // Caller ID имя абонента
    }


Структура с полной информацией о зарегистрированном PJSIP endpoint'е.

**Поля:**
- Extension - номер extension (внутренний номер)
- Contacts - полный SIP Contact URI в формате `sip:extension@ip:port;параметры`
- IP - IP адрес телефона (извлекается автоматически из Contacts)
- Name` - имя абонента из Caller ID

**Пример данных:**

    ExtensionData{
        Extension: "100",
        Contacts:  "sip:100@10.0.0.50:5060;transport=UDP",
        IP:        "10.0.0.50",
        Name:      "John Doe",
    }

### Работа с chan_sip
Для получения данных с Asterisk с chan_sip используйте GetSIPEndpointsIPtoDataMap  
Она возвращает то же самое за исключением того что поле Contacts пустое

### Тестирование
Можно выполнить интеграционный тест:  

    go test -v -args -pjsip <IP адрес Asterisk с PJSIP>\
    -pjsipu <Имя пользователя AMI> \
    -pjsipp <Пароль> \
    -chansip <IP адрес Asterisk с chan_sip> \
    -chansipu <Имя пользователя AMI> \
    -chansipp <Пароль>