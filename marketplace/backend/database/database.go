package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"marketplace/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string) {
	var err error
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к базе данных MarketplaceDB: %v", err)
	}

	log.Println("✅ Успешно подключено к базе Microsoft SQL Server: MarketplaceDB!")

	// Автомиграция всех сущностей
	err = DB.AutoMigrate(
		&models.User{},
		&models.Shop{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.SubOrder{},
		&models.OrderItem{},
		&models.Review{},
	)
	if err != nil {
		log.Printf("⚠️ Заметка по миграции: %v", err)
	}

	seedData()
}

func seedData() {
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("ℹ️ Данные в базе MarketplaceDB уже инициализированы.")
		return
	}

	log.Println("🌱 Заполнение начальными демонстрационными данными (Пользователи, Магазины, Категории, Товары, Заказы, Отзывы)...")

	defaultPassword := "123456"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	passStr := string(hashedPassword)

	buyer := models.User{
		Name:         "Дилшод Раҳимов",
		Email:        "buyer@market.tj",
		PasswordHash: passStr,
		Role:         "buyer",
		Phone:        "+992 90 111 2233",
		AvatarURL:    "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150",
	}
	DB.Create(&buyer)

	seller1User := models.User{
		Name:         "Парвиз Алиев (TechStore)",
		Email:        "seller1@market.tj",
		PasswordHash: passStr,
		Role:         "seller",
		Phone:        "+992 91 222 3344",
		AvatarURL:    "https://images.unsplash.com/photo-1570295999919-56ceb5ecca61?w=150",
	}
	DB.Create(&seller1User)

	seller2User := models.User{
		Name:         "Нигина Саидова (ModnoStyle)",
		Email:        "seller2@market.tj",
		PasswordHash: passStr,
		Role:         "seller",
		Phone:        "+992 92 333 4455",
		AvatarURL:    "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150",
	}
	DB.Create(&seller2User)

	moderatorUser := models.User{
		Name:         "Фаридун Умаров (Модератор)",
		Email:        "moderator@market.tj",
		PasswordHash: passStr,
		Role:         "moderator",
		Phone:        "+992 93 444 5566",
		AvatarURL:    "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150",
	}
	DB.Create(&moderatorUser)

	shop1 := models.Shop{
		UserID:      seller1User.ID,
		Name:        "TechStore",
		Slug:        "techstore",
		Description: "Мағозаи расмии техникаи рақамӣ, смартфонҳо, компютерҳо ва лавозимот бо кафолати аслӣ.",
		LogoURL:     "https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=200",
		BannerURL:   "https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200",
		Phone:       "+992 91 222 3344",
		Address:     "ш. Душанбе, хиёбони Рӯдакӣ, 45",
		Rating:      4.9,
		IsVerified:  true,
	}
	DB.Create(&shop1)

	shop2 := models.Shop{
		UserID:      seller2User.ID,
		Name:        "ModnoStyle",
		Slug:        "modnostyle",
		Description: "Либосҳои замонавӣ, пойафзол ва лавозимоти муосир барои мардон ва занон бо нархи дастрас.",
		LogoURL:     "https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=200",
		BannerURL:   "https://images.unsplash.com/photo-1441984904996-e0b6ba687e04?w=1200",
		Phone:       "+992 92 333 4455",
		Address:     "ш. Хуҷанд, кӯчаи И. Сомонӣ, 12",
		Rating:      4.8,
		IsVerified:  true,
	}
	DB.Create(&shop2)

	categories := []models.Category{
		{
			Name:        "Электроника",
			Slug:        "electronics",
			Icon:        "Laptop",
			ImageURL:    "https://images.unsplash.com/photo-1498049794561-7780e7231661?w=400",
			Description: "Смартфонҳо, ноутбукҳо, гӯшмонакҳо ва техникаи маишӣ",
			SortOrder:   1,
		},
		{
			Name:        "Либос ва пойафзол",
			Slug:        "clothing",
			Icon:        "Shirt",
			ImageURL:    "https://images.unsplash.com/photo-1489987707025-afc232f7ea0f?w=400",
			Description: "Либоси мардона, занона ва бачагонаи зебо",
			SortOrder:   2,
		},
		{
			Name:        "Хона ва ошхона",
			Slug:        "home-kitchen",
			Icon:        "Home",
			ImageURL:    "https://images.unsplash.com/photo-1513694203232-719a280e022f?w=400",
			Description: "Зарфҳо, нассоҷӣ, рӯзгор ва ороиши хона",
			SortOrder:   3,
		},
		{
			Name:        "Зебоӣ ва саломатӣ",
			Slug:        "beauty",
			Icon:        "Sparkles",
			ImageURL:    "https://images.unsplash.com/photo-1522335789203-aabd1fc54bc9?w=400",
			Description: "Пардоз, атр, маҳсулоти нигоҳубини мӯй ва пӯст",
			SortOrder:   4,
		},
		{
			Name:        "Варзиш ва сайёҳӣ",
			Slug:        "sports",
			Icon:        "Activity",
			ImageURL:    "https://images.unsplash.com/photo-1517838277536-f5f99be501cd?w=400",
			Description: "Анҷомҳои варзишӣ, хаймаҳо ва лавозимоти фитнес",
			SortOrder:   5,
		},
		{
			Name:        "Китобҳо",
			Slug:        "books",
			Icon:        "BookOpen",
			ImageURL:    "https://images.unsplash.com/photo-1495446815901-a7297e633e8d?w=400",
			Description: "Адабиёти бадеӣ, илмӣ, дарсӣ ва бачагона",
			SortOrder:   6,
		},
	}
	for i := range categories {
		DB.Create(&categories[i])
	}

	pOld := func(v float64) *float64 { return &v }

	products := []models.Product{

		{
			ShopID:        shop1.ID,
			CategoryID:    categories[0].ID,
			Name:          "Ноутбук Asus ZenBook 14 OLED Ultra 7",
			Slug:          "asus-zenbook-14-oled",
			Description:   "Ноутбуки ултрабуки борик бо экрани 2.8K OLED, протсессори пуриқтидори Intel Core Ultra 7, 32GB RAM ва 1TB SSD. Идеалӣ барои кор, барномасозӣ ва дизайн.",
			Price:         14500.00,
			OldPrice:      pOld(16200.00),
			StockQuantity: 15,
			ImageURL:      "https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=600",
			AdditionalImages: mustJSONString([]string{
				"https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=600",
				"https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=600",
			}),
			Status:       "approved",
			Rating:       4.9,
			ReviewsCount: 8,
			SalesCount:   24,
		},
		{
			ShopID:        shop1.ID,
			CategoryID:    categories[0].ID,
			Name:          "Смартфон Apple iPhone 15 Pro 256GB Titanium",
			Slug:          "iphone-15-pro-256gb",
			Description:   "Смартфони флагманӣ бо корпуси титанӣ, протсессори A17 Pro ва камераи пешрафтаи 48MP бо зуми 3x оптикӣ. Ранги табиии титан.",
			Price:         13800.00,
			OldPrice:      pOld(14900.00),
			StockQuantity: 12,
			ImageURL:      "https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?w=600",
			Status:        "approved",
			Rating:        5.0,
			ReviewsCount:  14,
			SalesCount:    35,
		},
		{
			ShopID:        shop1.ID,
			CategoryID:    categories[0].ID,
			Name:          "Гӯшмонакҳои бесими Sony WH-1000XM5",
			Slug:          "sony-wh-1000xm5",
			Description:   "Гӯшмонакҳои беҳтарини ҷаҳон бо технологияи пешрафтаи пахшкунии садо (Noise Cancelling), батарея то 30 соат ва садои Hi-Res Audio.",
			Price:         3450.00,
			OldPrice:      pOld(3900.00),
			StockQuantity: 28,
			ImageURL:      "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600",
			Status:        "approved",
			Rating:        4.8,
			ReviewsCount:  9,
			SalesCount:    19,
		},
		{
			ShopID:        shop1.ID,
			CategoryID:    categories[0].ID,
			Name:          "Смарт-соат Samsung Galaxy Watch 6 Classic",
			Slug:          "samsung-galaxy-watch-6",
			Description:   "Соати ҳушманд бо безели чархзананда, мониторинги хоб, санҷиши набз ва ЭКГ. Ба об тобовар 5ATM.",
			Price:         2890.00,
			OldPrice:      pOld(3300.00),
			StockQuantity: 20,
			ImageURL:      "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600",
			Status:        "approved",
			Rating:        4.7,
			ReviewsCount:  6,
			SalesCount:    15,
		},

		{
			ShopID:        shop1.ID,
			CategoryID:    categories[0].ID,
			Name:          "Планшет Apple iPad Air M2 11-inch 128GB",
			Slug:          "ipad-air-m2-11",
			Description:   "Планшети навтарин бо чипи пуриқтидори Apple M2, экрани Liquid Retina ва дастгирии Apple Pencil Pro.",
			Price:         7900.00,
			OldPrice:      pOld(8500.00),
			StockQuantity: 10,
			ImageURL:      "https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=600",
			Status:        "pending",
			Rating:        0,
			ReviewsCount:  0,
			SalesCount:    0,
		},

		{
			ShopID:        shop2.ID,
			CategoryID:    categories[1].ID,
			Name:          "Курткаи чармии мардона Classic Biker",
			Slug:          "leather-jacket-biker",
			Description:   "Куртка аз чарми табиӣ бо фурнитураи мустаҳкам. Дизайни замонавӣ ва бароҳат барои фаслҳои баҳору тирамоҳ.",
			Price:         1250.00,
			OldPrice:      pOld(1600.00),
			StockQuantity: 22,
			ImageURL:      "https://images.unsplash.com/photo-1520975916090-3105956dac38?w=600",
			Status:        "approved",
			Rating:        4.9,
			ReviewsCount:  11,
			SalesCount:    42,
		},
		{
			ShopID:        shop2.ID,
			CategoryID:    categories[1].ID,
			Name:          "Кросссовкаҳои варзишии Nike Air Max Pulse",
			Slug:          "nike-air-max-pulse",
			Description:   "Пойафзоли варзишии бароҳат барои давидан ва истифодаи ҳаррӯза бо болиштчаи ҳавоии нарм.",
			Price:         980.00,
			OldPrice:      pOld(1250.00),
			StockQuantity: 30,
			ImageURL:      "https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600",
			Status:        "approved",
			Rating:        5.0,
			ReviewsCount:  15,
			SalesCount:    58,
		},
		{
			ShopID:        shop2.ID,
			CategoryID:    categories[1].ID,
			Name:          "Либоси шевои занонаи шомгоҳӣ Emerald Silk",
			Slug:          "emerald-silk-dress",
			Description:   "Либоси абрешими нафис бо ранги изумрудӣ. Барои ҷашнвораҳо, тӯйҳо ва чорабиниҳои расмӣ.",
			Price:         850.00,
			OldPrice:      pOld(1100.00),
			StockQuantity: 18,
			ImageURL:      "https://images.unsplash.com/photo-1515372039744-b8f02a3ae446?w=600",
			Status:        "approved",
			Rating:        4.8,
			ReviewsCount:  7,
			SalesCount:    26,
		},
		{
			ShopID:        shop2.ID,
			CategoryID:    categories[1].ID,
			Name:          "Костюми классикии мардона Slim Fit Navy",
			Slug:          "mens-suit-slim-fit",
			Description:   "Костюми босифати 2-қисма (пиджак ва шим) аз пахта ва пашми табиӣ. Буриши дақиқ ва олии замонавӣ.",
			Price:         1950.00,
			OldPrice:      pOld(2400.00),
			StockQuantity: 14,
			ImageURL:      "https://images.unsplash.com/photo-1594938298603-c8148c4dae35?w=600",
			Status:        "approved",
			Rating:        4.9,
			ReviewsCount:  5,
			SalesCount:    18,
		},

		{
			ShopID:        shop2.ID,
			CategoryID:    categories[1].ID,
			Name:          "Пальтои гарми зимистонаи кашемирӣ",
			Slug:          "cashmere-winter-coat",
			Description:   "Пальтои дароз бо матои кашемири 100% бо ҷайбҳои калон ва камарбанди чармӣ.",
			Price:         1600.00,
			OldPrice:      pOld(2100.00),
			StockQuantity: 8,
			ImageURL:      "https://images.unsplash.com/photo-1539533018447-63fcce667823?w=600",
			Status:        "pending",
			Rating:        0,
			ReviewsCount:  0,
			SalesCount:    0,
		},

		{
			ShopID:        shop2.ID,
			CategoryID:    categories[2].ID,
			Name:          "Маҷмӯи зарфҳои қаҳвапазии керамикӣ (6 нафарӣ)",
			Slug:          "ceramic-coffee-set",
			Description:   "Маҷмӯи қаҳванушӣ аз сафоли табиӣ бо ороишоти дастӣ ва табақчаҳои чӯбӣ.",
			Price:         420.00,
			OldPrice:      pOld(550.00),
			StockQuantity: 25,
			ImageURL:      "https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=600",
			Status:        "approved",
			Rating:        4.7,
			ReviewsCount:  4,
			SalesCount:    16,
		},

		{
			ShopID:        shop1.ID,
			CategoryID:    categories[4].ID,
			Name:          "Велосипеди кӯҳии касбии Trek Marlin 7",
			Slug:          "trek-marlin-7-bike",
			Description:   "Велосипеди пуриқтидор бо тормозҳои гидравликии дискӣ, амортизатори танзимшаванда ва чархҳои 29 дюйма.",
			Price:         5800.00,
			OldPrice:      pOld(6500.00),
			StockQuantity: 7,
			ImageURL:      "https://images.unsplash.com/photo-1485965120184-e220f721d03e?w=600",
			Status:        "approved",
			Rating:        5.0,
			ReviewsCount:  6,
			SalesCount:    12,
		},
	}

	for i := range products {
		DB.Create(&products[i])
	}

	reviews := []models.Review{
		{
			ProductID: products[0].ID,
			UserID:    buyer.ID,
			Rating:    5,
			Comment:   "Ноутбуки бисёр олиҷаноб! Экрани OLED чашмро хаста намекунад, тез кор мекунад ва батареяаш тамоми рӯз мерасад. Ба ҳама тавсия медиҳам!",
			CreatedAt: time.Now().Add(-10 * 24 * time.Hour),
		},
		{
			ProductID: products[1].ID,
			UserID:    buyer.ID,
			Rating:    5,
			Comment:   "Сифати аслӣ! Телефон дар ҳақиқат аслии Apple аст ва расонидан хеле зуд буд. Ташаккур ба TechStore!",
			CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ProductID: products[5].ID,
			UserID:    buyer.ID,
			Rating:    5,
			Comment:   "Куртка хеле хуб дӯхта шудааст. Андозааш дақиқ рост омад, чармаш нарм ва бӯи хуш дорад.",
			CreatedAt: time.Now().Add(-15 * 24 * time.Hour),
		},
	}
	for i := range reviews {
		DB.Create(&reviews[i])
	}

	order1Date := time.Now().Add(-3 * 24 * time.Hour)
	order1 := models.Order{
		OrderNumber:     "ORD-2026-9001",
		UserID:          buyer.ID,
		TotalAmount:     4430.00,
		Status:          "processing",
		RecipientName:   "Дилшод Раҳимов",
		RecipientPhone:  "+992 90 111 2233",
		ShippingAddress: "ш. Душанбе, н. Сино, кӯчаи Борбад 48, хонаи 12",
		PaymentMethod:   "cash_on_delivery",
		Notes:           "Илтимос пеш аз омадан занг занед",
		CreatedAt:       order1Date,
		UpdatedAt:       order1Date,
	}
	DB.Create(&order1)

	subOrder1Shop1 := models.SubOrder{
		OrderID:          order1.ID,
		ShopID:           shop1.ID,
		SubtotalAmount:   3450.00,
		Status:           "processing",
		ShippingTracking: "TRACK-TJ-TECH-101",
		CreatedAt:        order1Date,
		UpdatedAt:        order1Date,
	}
	DB.Create(&subOrder1Shop1)

	item1Shop1 := models.OrderItem{
		OrderID:      order1.ID,
		SubOrderID:   subOrder1Shop1.ID,
		ShopID:       shop1.ID,
		ProductID:    products[2].ID,
		ProductName:  products[2].Name,
		ProductPrice: products[2].Price,
		Quantity:     1,
		TotalPrice:   3450.00,
		ImageURL:     products[2].ImageURL,
		CreatedAt:    order1Date,
	}
	DB.Create(&item1Shop1)

	subOrder1Shop2 := models.SubOrder{
		OrderID:          order1.ID,
		ShopID:           shop2.ID,
		SubtotalAmount:   980.00,
		Status:           "shipped",
		ShippingTracking: "TRACK-TJ-MOD-202",
		CreatedAt:        order1Date,
		UpdatedAt:        order1Date,
	}
	DB.Create(&subOrder1Shop2)

	item1Shop2 := models.OrderItem{
		OrderID:      order1.ID,
		SubOrderID:   subOrder1Shop2.ID,
		ShopID:       shop2.ID,
		ProductID:    products[6].ID,
		ProductName:  products[6].Name,
		ProductPrice: products[6].Price,
		Quantity:     1,
		TotalPrice:   980.00,
		ImageURL:     products[6].ImageURL,
		CreatedAt:    order1Date,
	}
	DB.Create(&item1Shop2)

	seedMonthlyPastOrders(buyer.ID, shop1.ID, shop2.ID, products)

	log.Println("✅ Начальные данные базы MarketplaceDB успешно созданы!")
}

func seedMonthlyPastOrders(buyerID, shop1ID, shop2ID int64, products []models.Product) {
	now := time.Now()

	lastMonthDate := now.AddDate(0, -1, -5)
	orderPast1 := models.Order{
		OrderNumber:     fmt.Sprintf("ORD-%d-7042", lastMonthDate.Year()),
		UserID:          buyerID,
		TotalAmount:     15750.00,
		Status:          "completed",
		RecipientName:   "Дилшод Раҳимов",
		RecipientPhone:  "+992 90 111 2233",
		ShippingAddress: "ш. Душанбе, н. И. Сомонӣ",
		PaymentMethod:   "card",
		CreatedAt:       lastMonthDate,
		UpdatedAt:       lastMonthDate,
	}
	DB.Create(&orderPast1)

	subPast1 := models.SubOrder{
		OrderID:          orderPast1.ID,
		ShopID:           shop1ID,
		SubtotalAmount:   14500.00,
		Status:           "delivered",
		ShippingTracking: "TRACK-OCT-01",
		CreatedAt:        lastMonthDate,
		UpdatedAt:        lastMonthDate,
	}
	DB.Create(&subPast1)

	DB.Create(&models.OrderItem{
		OrderID:      orderPast1.ID,
		SubOrderID:   subPast1.ID,
		ShopID:       shop1ID,
		ProductID:    products[0].ID,
		ProductName:  products[0].Name,
		ProductPrice: 14500.00,
		Quantity:     1,
		TotalPrice:   14500.00,
		ImageURL:     products[0].ImageURL,
		CreatedAt:    lastMonthDate,
	})

	subPast2 := models.SubOrder{
		OrderID:          orderPast1.ID,
		ShopID:           shop2ID,
		SubtotalAmount:   1250.00,
		Status:           "delivered",
		ShippingTracking: "TRACK-OCT-02",
		CreatedAt:        lastMonthDate,
		UpdatedAt:        lastMonthDate,
	}
	DB.Create(&subPast2)

	DB.Create(&models.OrderItem{
		OrderID:      orderPast1.ID,
		SubOrderID:   subPast2.ID,
		ShopID:       shop2ID,
		ProductID:    products[5].ID,
		ProductName:  products[5].Name,
		ProductPrice: 1250.00,
		Quantity:     1,
		TotalPrice:   1250.00,
		ImageURL:     products[5].ImageURL,
		CreatedAt:    lastMonthDate,
	})

	twoMonthsDate := now.AddDate(0, -2, -10)
	orderPast2 := models.Order{
		OrderNumber:     fmt.Sprintf("ORD-%d-6021", twoMonthsDate.Year()),
		UserID:          buyerID,
		TotalAmount:     16690.00,
		Status:          "completed",
		RecipientName:   "Дилшод Раҳимов",
		RecipientPhone:  "+992 90 111 2233",
		ShippingAddress: "ш. Душанбе",
		PaymentMethod:   "card",
		CreatedAt:       twoMonthsDate,
		UpdatedAt:       twoMonthsDate,
	}
	DB.Create(&orderPast2)

	subPast3 := models.SubOrder{
		OrderID:          orderPast2.ID,
		ShopID:           shop1ID,
		SubtotalAmount:   13800.00,
		Status:           "delivered",
		ShippingTracking: "TRACK-SEP-01",
		CreatedAt:        twoMonthsDate,
		UpdatedAt:        twoMonthsDate,
	}
	DB.Create(&subPast3)

	DB.Create(&models.OrderItem{
		OrderID:      orderPast2.ID,
		SubOrderID:   subPast3.ID,
		ShopID:       shop1ID,
		ProductID:    products[1].ID,
		ProductName:  products[1].Name,
		ProductPrice: 13800.00,
		Quantity:     1,
		TotalPrice:   13800.00,
		ImageURL:     products[1].ImageURL,
		CreatedAt:    twoMonthsDate,
	})

	subPast4 := models.SubOrder{
		OrderID:          orderPast2.ID,
		ShopID:           shop2ID,
		SubtotalAmount:   2890.00,
		Status:           "delivered",
		ShippingTracking: "TRACK-SEP-02",
		CreatedAt:        twoMonthsDate,
		UpdatedAt:        twoMonthsDate,
	}
	DB.Create(&subPast4)

	DB.Create(&models.OrderItem{
		OrderID:      orderPast2.ID,
		SubOrderID:   subPast4.ID,
		ShopID:       shop2ID,
		ProductID:    products[8].ID,
		ProductName:  products[8].Name,
		ProductPrice: 1950.00,
		Quantity:     1,
		TotalPrice:   1950.00,
		ImageURL:     products[8].ImageURL,
		CreatedAt:    twoMonthsDate,
	})
}

func mustJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
