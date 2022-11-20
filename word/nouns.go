package word

// TODO: 补充更多表示具体物品的名词
// 这个列表中的单词允许重复
var nouns = []string{
	/* 第一批整理的单词 */
	// 家具及家居用品
	"bed",        // 床
	"crib",       // 婴儿床
	"nightstand", // 床头柜
	"vanity",     // 梳妆台
	"dresser",    // 抽屉柜
	"wardrobe",   // 衣柜
	"bookcase",   // 书柜
	"organizer",  // 收纳柜
	"bench",      // 长板凳
	// 日用品
	"battery",    // 电池
	"candle",     // 蜡烛
	"cotton",     // 棉
	"envelopes",  // 信封
	"lighter",    // 打火机
	"matches",    // 火柴
	"needle",     // 针
	"scissors",   // 剪刀
	"sellotape",  // 胶带
	"stamps",     // 邮票
	"pen",        // 钢笔
	"pencil",     // 铅笔
	"tissues",    // 纸巾
	"toothpaste", // 牙膏
	"soap",       // 肥皂
	// 教室用品
	"map",        // 地图
	"chalkboard", // 黑板
	"globe",      // 地球；地球仪；球体
	"eraser",     // 橡皮
	"chalk",      // 粉笔
	"desk",       // 课桌
	"chair",      // 椅子
	"textbook",   // 教科书，课本
	// 常见衣服
	"tie",        // 领带
	"scarf",      // 围巾
	"gloves",     // 手套
	"trousers",   // 男士长裤
	"jacket",     // 夹克衫
	"shirt",      // 衬衫
	"skirt",      // 短裙子
	"dress",      // 连衣裙
	"pants",      // 女士长裤
	"bikini",     // 比基尼
	"swimsuit",   // 泳衣
	"jeans",      // 牛仔裤
	"socks",      // 袜子
	"shoes",      // 鞋子
	"sweater",    // 毛衣
	"coat",       // 上衣
	"raincoat",   // 雨衣
	"shorts",     // 短裤
	"sneakers",   // 网球鞋
	"slippers",   // 拖鞋
	"sandals",    // 凉鞋
	"boots",      // 靴子
	"hat",        // (有沿的)帽子
	"cap",        // 便帽
	"sunglasses", // 太阳镜
	"goggles ",   // 泳镜
	// 身体器官名词
	"head",     // 头
	"hair",     // 头发
	"skull",    // 颅骨，头盖骨
	"brain",    // 脑
	"forehead", // 额
	"eyebrow",  // 眉毛
	"eyelash",  // 眼睫毛
	"eye",      // 眼睛
	"nose",     // 鼻子
	"mouth",    // 口
	"lip",      // 嘴唇
	"tooth",    // 牙齿
	"tongue",   // 舌
	"ear",      // 耳朵
	"cheek",    // 面颊
	"chin",     // 下巴
	"neck",     // 脖子
	"throat",   // 喉咙，咽喉
	"shoulder", // 肩
	"arm",      // 手臂
	"elbow",    // 肘
	"wrist",    // 手腕
	// "hand",       // 手
	// "palm",       // 手掌
	"finger",     // 手指
	"nail",       // 指甲
	"thumb",      // 大拇指
	"forefinger", // 食指
	"fist",       // 拳头
	// "knuckle",    // 指关节
	"backbone",   // 脊骨，脊柱
	"collarbone", // 锁骨
	// "chest",      // 胸部
	"breastbone", // 胸骨
	"heart",      // 心脏
	"lung",       // 肺
	"abdomen",    // 腹部
	"belly",      // 肚子
	"stomach",    // 胃
	"back",       // 背
	"waist",      // 腰
	"pelvis",     // 骨盆
	"buttocks",   // 屁股
	"leg",        // 腿
	"thigh",      // 大腿
	"knee",       // 膝盖
	// "shank",      // 小腿
	"ankle", // 脚踝
	"foot",  // 脚
	// "instep", // 脚背
	// "heel", // 脚后跟
	// "sole", // 脚底
	// "arch", // 脚掌心
	"toes", // 脚趾

	/* 来自 https://new.qq.com/rain/a/20201207A0HZS600.html，有删改*/
	// 房屋外部
	"doghouse",  // 犬屋
	"curtain",   // 窗帘
	"garage",    // 车房,车库
	"porch",     // 入口处
	"driveway",  // 车库通向马路的空地
	"mailbox",   // 信箱
	"dormer",    // 屋顶窗
	"skylight",  // 天窗
	"chimney",   // 烟囱
	"balcony",   // 阳台
	"shutter",   // 百叶窗
	"lawn",      // 草坪,草地
	"shrubs",    // 灌木
	"sprinkler", // 自动撒水器
	// 房屋内部
	"window",     // 窗户
	"television", // 电视机
	"console",    // 主控台 控制台
	"chair",      // 椅子
	"floor",      // 地板,地面
	"carpet",     // 地毯
	"clock",      // 钟
	"calendar",   // 日历
	"door",       // 门
	"bookcase",   // 书柜，书橱
	"couch",      // 沙发
	"lamp",       // 灯
	"wall",       // 墙
	// 生活用品
	"lighter",   // 打火机
	"matches",   // 火柴
	"ashtray",   // 烟灰缸
	"cigarette", // 香烟
	"armchair",  // 扶手椅
	"vase",      // 花瓶
	"telephone", // 电话机
	"recliner",  // 卧椅
	// 厨房
	"refrigerator", // 冰箱
	"counter",      // 柜台
	"sink",         // 洗涤槽,水槽
	"wok",          // 铁锅(带把的中国炒菜锅)
	"pan",          // 平底锅
	"ladle",        // 勺子;长柄勺
	"ventilator",   // 通风机;换气扇
	"apron",        // 围裙;工作裙
	"cupboard",     // 食橱;碗柜
	"oven",         // 炉,灶
	"cabinets",     // 橱柜
	"dustpan",      // 簸箕
	"broom",        // 扫帚
	"mop",          // 拖把
	"blender",      // 搅拌机，捣碎机
	"toaster",      // 烤面包器;烤炉,烤箱
	"knife",        // 刀,小刀;菜刀;
	"microwave",    // 微波炉
	// 浴室
	"showerhead", // 喷头
	"faucet",     // 龙头,旋塞
	"bathtub",    // 浴缸
	"drain",      // 排水管,下水道
	"tile",       // 瓦;瓷砖;墙砖;地砖
	"mirror",     // 镜子
	"sink",       // 洗涤槽,水槽
	"washcloth",  // 毛巾
	"cabinet",    // 橱,柜
	"rug",        // (铺于室内部分地面上的)小地毯;毛皮地毯
	"toilet",     // (有冲洗式马桶的)厕所,洗手间,盥洗室
	"reservoir",  // 蓄水库;贮水池(或槽);贮存器
	// 飞机
	"captain",    // 机长
	"copilot",    // 副驾驶员
	"steward",    // 男空服员
	"fuselage",   // (飞机的)机身
	"lavatory",   // 盥洗室
	"stewardess", // 空姐
	// 机场
	"landing",   // 降落
	"passenger", // 乘客
	"luggage",   // 行李
	"takeoff",   // 起飞
	"airplane",  // 飞机
	"airlines",  // (飞机的)航线
	"customs",   // 进口税, 海关
	// 教室
	"map",        // 地图
	"chalkboard", // 黑板
	"organ",      // 管风琴
	"globe",      // 地球仪
	"eraser",     // 黑板擦
	"chalk",      // 粉笔
	"platform",   // 讲台
	"desk",       // 书桌
	"chair",      // 椅子
	"textbook",   // 教材
	// 商场
	"elevator",  // 电梯
	"counter",   // 柜台
	"umbrella",  // 雨伞
	"escalator", // 自动楼梯
	"desserts",  // 甜点心;餐后甜点
	"tie",       // 领带
	"coat",      // 外套
	"socks",     // 短袜
	"briefcase", // 公事包
	"radio",     // 收音机
	"gloves",    // 手套
	"shoes",     // 鞋子
	"basement",  // 地下室
	"stairs",    // 楼梯
	"sweater",   // 毛线衣
	"skirt",     // 裙子;衬裙
	// 小卖部
	"poster",    // 海报
	"candy",     // 糖果
	"pretzel",   // 双圈饼干
	"cookie",    // 曲奇饼
	"gum",       // 口香糖
	"mustard",   // 芥末
	"popsicle",  // 棒冰
	"drumstick", // 鸡腿
	"popcorn",   // 爆米花
	"hamburger", // 汉堡
	"ketchup",   // 番茄酱
	"fries",     // 薯条
	// 自助洗衣店
	"dryer",       // 烘干机
	"screen",      // 筛
	"washboard",   // 搓衣板
	"detergent",   // 洗涤剂
	"starch",      // 浆粉
	"washtub",     // 洗衣盆
	"bleach",      // 漂白剂
	"hamper",      // 篮子
	"lint",        // 软麻布
	"clothesline", // 晾衣绳
	"stain",       // 污点
	// 园艺
	"shed",       // 小屋
	"shrubs",     // 灌木丛
	"cuttings",   // 插条
	"rake",       // 耙
	"hoe",        // 锄
	"spade",      // 铲
	"seeds",      // 种子
	"bulbs",      // 鳞茎
	"shears",     // 大剪刀
	"hose",       // 软管
	"sprinkler",  // 洒水器
	"fertilizer", // 肥料
	"dung",       // 粪
	"weeds",      // 杂草
	// 农场
	"field",     // 田地
	"haystack",  // 干草堆
	"wagon",     // 四轮车
	"hay",       // 干草
	"rope",      // 绳子
	"fence",     // 篱笆
	"hoe",       // 锄
	"rake",      // 耙
	"horseshoe", // 马掌
	"shovel",    // 铲
	"corn",      // 谷物
	"pitchfork", // 干草叉
	"horse",     // 马
	"trough",    // 饲料槽
	"barn",      // 谷仓
	"tractor",   // 拖拉机
	"wheat",     // 小麦
	// 游乐园
	"games",    // 游戏
	"vomit",    // 呕吐
	"line",     // 队伍
	"carousel", // 旋转木马
	"tracks",   // 轨道
	"Seat",     // 座位
	"Car",      // 车厢
	// 牙科诊所
	"Sink",       // 水槽
	"towel",      // 毛巾
	"Mirror",     // 镜子
	"Tools",      // 工具
	"Tray",       // 托盘
	"Drill",      // 钻子
	"Toothbrush", // 牙刷
	"Dentures",   // 假牙
	"Mold",       // 模子
	"Teeth",      // 牙齿
	// 复活节
	"basket",       // 篮子
	"egg",          // 蛋
	"chocolate",    // 巧克力
	"handle",       // 手把
	"dye",          // 染料
	"marshmallows", // 棉花软糖
	"paint",        // 油漆
	"bunny",        // 兔子
	"whiskers",     // 胡须
	"ears",         // 耳朵
	"paws",         // 爪子
	// 化妆品
	"moisturizer", // 润肤乳
	"blush",       // 腮红
	"brushes",     // 刷子
	"tweezers",    // 镊子
	"lipstick",    // 口红
	"foundation",  // 粉底霜
	"eyeliner",    // 眼线笔
	"mascara",     // 睫毛膏
	// 发廊
	"sink",        // 水槽
	"mirror",      // 镜子
	"comb",        // 梳子
	"brush",       // 刷子
	"scissors",    // 剪子
	"shampoo",     // 洗发水
	"mousse",      // 摩斯
	"razor",       // 剃刀
	"towel",       // 毛巾
	"conditioner", // 护发素
	// 游泳池
	"Lifeguard", // 救生员
	"Umbrella",  // 伞
	"bikini",    // 比基尼
	"net",       // 网
	"lockers",   // 储物柜
	"earplugs",  // 耳塞
	"towel",     // 毛巾
	"goggles",   // 泳镜
	"sunscreen", // 防晒霜
	"ball",      // 球
	"kickboard", // 踢水板
	"ladder",    // 梯子
	// 衣橱
	"blazer",     // 小西装
	"suspenders", // 吊裤带
	"suit",       // 套装
	"tie",        // 领带
	"pants",      // 裤子
	"vest",       // 背心
	"raincoat",   // 雨衣
	"blouse",     // 女士衬衣
	"scarf",      // 围巾
	"purse",      // 女手提袋
	"skirt",      // 短裙子
	"coat",       // 大衣
	"hanger",     // 衣架
	"dress",      // 连衣裙
	"Shorts",     // 短裤
	"Socks",      // 袜子
	"Belt",       // 皮带
	"Jeans",      // 牛仔裤
	"Sweater",    // 毛衣
	"Drawer",     // 抽屉
	// 便当
	"thermos",  // 热水瓶
	"straw",    // 吸管
	"fork",     // 叉
	"knife",    // 刀
	"spoon",    // 勺
	"lid",      // 盖子
	"salad",    // 沙拉
	"latch",    // 闩
	"yogurt",   // 乳酪
	"hinge",    // 脚链
	"sandwich", // 三明治
	"pepper",   // 胡椒
	"napkin",   // 餐巾
	"Pepper",   // 胡椒
	// 托儿所
	"pacifier", // 奶嘴
	"bib",      // 围嘴
	"blanket",  // 毯子
	"rattle",   // 摇响器
	"diaper",   // 尿布
	"crib",     // 婴儿床
	"playpen",  // 游戏围栏
	"stroller", // 婴儿车
	"booties",  // 婴儿袜
	// 玩具箱
	"doll",    // 娃娃
	"slinky",  // 弹簧玩具
	"blocks",  // 积木
	"legos",   // 乐高积木
	"robot",   // 机器人
	"crayons", // 蜡笔
	"Barbie",  // 芭比娃娃
	"Car",     // 车
	"Ball",    // 球
	// 游艇
	"deck",     // 夹板
	"cocktail", // 鸡尾酒
	"sandals",  // 拖鞋
	"seagull",  // 海鸥
	"bikini",   // 比基尼
	"anchor",   // 锚
	"rope",     // 绳子
	"captain",  // 船长
	"oar",      // 浆
	"porthole", // 舷窗
	// 水果超市
	"peach",      // 桃子
	"grapes",     // 葡萄
	"basket",     // 篮子
	"papaya",     // 木瓜
	"guava",      // 番石榴
	"lime",       // 青柠
	"lemon",      // 柠檬
	"seeds",      // 种子
	"nectarine",  // 油桃
	"scale",      // 秤
	"grapefruit", // 葡萄柚
	"plum",       // 李子
	"mango",      // 芒果
	"taro",       // 芋头
	// 管弦乐团
	"bass",      // （低音部）
	"harp",      // 竖琴
	"conductor", // 指挥
	"cello",     // 大提琴
	"bow",       // 弓
	"baton",     // 指挥棒
	"drums",     // 鼓
	"clarinet",  // 单簧管
	"violin",    // 小提琴
	"viola",     // 中提琴
	"saxophone", // 萨克斯
	"flute",     // 长笛
	"trombone",  // 长号
	"trumpet",   // 小号
	// 潜水
	"fins",       // 蛙鞋
	"mask",       // 面镜
	"compass",    // 罗盘
	"snorkel",    // 呼吸管
	"gauge",      // 潜水计量器
	"regulator",  // 调节器
	"coral",      // 珊瑚
	"booties",    // 潜水靴
	"boat",       // 船
	"glove",      // 潜水手套
	"flashlight", // 手电筒
	// 工具箱
	"screws",       // 螺丝
	"screwdriver",  // 螺丝刀
	"bolts",        // 螺栓
	"saw",          // 锯子
	"washers",      // 垫圈
	"nuts",         // 螺母
	"hammer",       // 锤子
	"wrench",       // 扳手
	"nails",        // 钉子
	"toolbox",      // 工具箱
	"level",        // 水平
	"sledgehammer", // 大锤
	"drill",        // 钻孔机
	"tape",         // 胶带
	// 诊所
	"thermometer", // 温度计
	"syringe",     // 注射器
	"eardrops",    // 耳药水
	"medicine",    // 药
	"pills",       // 药片
	"antibiotics", // 抗生素
	"painkillers", // 止痛片
	"beaker",      // 烧杯
	"stethoscope", // 听诊器
	"scale",       // 秤
	"swab",        // 棉签
	// 露天咖啡座
	"chair",     // 椅子
	"ashtray",   // 烟灰缸
	"muffin",    // 松饼
	"stirrer",   // 搅拌棒
	"top",       // 杯盖
	"chocolate", // 巧克力
	"menu",      // 菜单
	"napkin",    // 餐巾
	"cream",     // 泡沫
	"umbrella",  // 伞
	"cup",       // 杯子
	"saucer",    // 杯垫
	"table",     // 桌子
	// 捕鱼
	"Bait",   // 诱饵
	"Worm",   // 虫
	"waders", // 钓鱼裤
	"reel",   // 转轮
	"hook",   // 钩子
	"net",    // 网
	"lure",   // 诱饵
	"fly",    // 人工拟饵
	// 加油站
	"squeegee",  // 清洁器
	"attendant", // 服务员
	"restroom",  // 卫生家
	"tire",      // 轮胎
	// 感恩节
	"pasta",    // 意大利面
	"knife",    // 刀子
	"sausages", // 香肠
	"turkey",   // 火鸡
	"muffin",   // 松饼
	"biscuits", // 饼干
	"butter",   // 黄油
	"soup",     // 汤
	"gravy",    // 肉汁
	"stuffing", // 填充料
	"pudding",  // 布丁
	"corn",     // 玉米
	"salad",    // 沙拉
	"squash",   // 西葫芦
	"yams",     // 山药
	// 自行车
	"helmet",     // 头盔
	"spoke",      // 辐条
	"lock",       // 锁
	"seat",       // 座子
	"gears",      // 齿轮
	"rim",        // 轮辋
	"chain",      // 链条
	"kickstand",  // 支架
	"derailleur", // 变速器
	"wheels",     // 轮胎
	"headset",    // 车头碗组
	"pedal",      // 踏板
	"brake",      // 刹车
	"reflector",  // 反光灯
	"basket",     // 车篮
	"shock",      // 避震条
}
