# 球线库增量资料与上线交接（第二批）

核对日期：2026-10-27。新增31款、82个毫米规格；本地 Docker 测试库总计12品牌、71款、184规格。按用户要求未备份。2026-10-27已同步生产库：12品牌、71款、184规格，全部启用。

数据文件：[expansion.v2.json](../data/tennis-strings/expansion.v2.json)。第一批仍为 [mainstream.v1.json](../data/tennis-strings/mainstream.v1.json)，保留原始核对信息。第二批 existingItemUpdates 覆盖第一批 X-ONE BIPHASE 未知材质。

## 来源与范围

以品牌官网商品页、区域官方商店和厂商目录为主。NRG2 的1.24毫米另由 Tennis Warehouse 商品页交叉确认。地区目录有售不等于中国渠道实时有库存，不能对外宣传为市面全量或销量排名。仅记录型号、已证实毫米线径和基本材料，不复制图片、评论、营销文案或实验性能数据。

型号的毫米规格为已确认子集。Gauge不跨品牌换算，商品编号不直接推断毫米值。大盘与小包装合并；不同配方或版本分开。材料只按证据归一；含添加剂不表示化学纯材料。

## 本批清单

|品牌|型号|线径 mm|来源与核对说明|
|---|---|---|---|
|Solinco|Revolution|1.15, 1.20, 1.25, 1.30|[来源1](https://solincosports.com/equipment/revolution/) |
|Solinco|Outlast|1.15, 1.20, 1.25, 1.30|[来源1](https://solincosports.com/equipment/outlast/) |
|Solinco|Barb Wire|1.20, 1.25, 1.30|[来源1](https://solincosports.com/equipment/barb-wire/) |
|Solinco|Tour Bite Diamond Rough|1.20, 1.25, 1.30|[来源1](https://solincosports.com/equipment/tour-bite-diamond-rough/) |
|Solinco|Vanquish|1.20, 1.30, 1.35|[来源1](https://solincosports.com/equipment/vanquish/) 保留官网明确列出的17/1.20规格，不按其他品牌Gauge换算。|
|Solinco|Mach-10|1.15, 1.20, 1.25, 1.30|[来源1](https://shop.solincosports.com/products/mach-10) |
|Signum Pro|X-perience|1.18, 1.24, 1.30|[来源1](https://www.signum-pro.com/html/index.php?id=86) |
|Signum Pro|Poly Plasma|1.18, 1.23, 1.28, 1.33|[来源1](https://www.signum-pro.com/html/index.php?id=5) |
|Signum Pro|Tornado|1.17, 1.23, 1.29|[来源1](https://www.signum-pro.com/html/index.php?id=42) |
|Signum Pro|Firestorm|1.20, 1.25, 1.30|[来源1](https://www.signum-pro.com/html/index.php?id=44) |
|Kirschbaum|Max Power|1.20, 1.25, 1.30|[来源1](https://kirschbaum-strings.de/produkte/produktdetail/max-power%3A126) |
|Kirschbaum|Pro Line|1.15, 1.20, 1.25, 1.30|[来源1](https://kirschbaum-strings.de/produkte/produktdetail/pro-line-ehemals-pro-line-ii%3A292) / [来源2](https://kirschbaum-strings.de/produkte/tennissaiten/co-polyester-saiten) 厂商明确Pro Line II改名为Pro Line且线材相同，旧名仅作别名，不重复建型号。材质由厂商Co-Polyester分类页确认。|
|Kirschbaum|Pro Line Rough|1.20, 1.25, 1.30|[来源1](https://kirschbaum-strings.de/produkte/produktdetail/pro-line-rough%3A139) |
|Kirschbaum|Evolution|1.20, 1.25, 1.30|[来源1](https://kirschbaum-strings.de/produkte/tennissaiten/co-polyester-saiten) |
|HEAD|Velocity MLT|1.25, 1.30|[来源1](https://www.head.com/en_US/product/velocity-mlt-4) 官网明确型号、结构与毫米线径；化学材质未在本次产品页明确，留空，不以结构推断。|
|HEAD|Lynx Touch|1.25, 1.30|[来源1](https://www.head.com/en_US/product/lynx-touch-281042) 官网明确型号、结构与毫米线径；化学材质未在本次产品页明确，留空，不以结构推断。|
|HEAD|Hawk Power|1.25, 1.30|[来源1](https://www.head.com/en_US/product/hawk-power-281022) 官网明确型号、结构与毫米线径；化学材质未在本次产品页明确，留空，不以结构推断。|
|GOSEN|G-TOUR 1|1.25, 1.30|[来源1](https://www.gosencanada.com/g-tour) 加拿大官方目录确认毫米规格；地区供应状态可能不同。|
|GOSEN|G-TOUR 2|1.25|[来源1](https://www.gosencanada.com/g-tour) 加拿大官方目录确认毫米规格；地区供应状态可能不同。|
|GOSEN|G-TOUR 3|1.18, 1.23, 1.28|[来源1](https://www.gosencanada.com/g-tour) 加拿大官方目录确认毫米规格；地区供应状态可能不同。|
|MSV|Focus HEX|1.10, 1.18, 1.23, 1.27, 1.33|[来源1](https://www.msv-tennis-shop.com/product/msv-focus-hex-tennis-string-12m-our-bestseller.63.html?language_code=en) / [来源2](https://www.msv-tennis-shop.com/documents/products/PDF/ERSA%20Testbericht%20MSV%20Focus-Hex%201%2C18%20Juli%202016%20k.pdf?ver=1739458771) 线径由官方商店确认；材质、结构由其托管的ERSA测试报告确认。|
|Dunlop|Explosive Bite|1.27, 1.32|[来源1](https://www.dunloptennis.net/storage/catalog2025.pdf) / [来源2](https://us.dunlopsports.com/tennis-strings.html) 日本官方2025目录明确毫米线径；为目录覆盖，不保证中国所有渠道当前有货。|
|Dunlop|Explosive Spin|1.25, 1.30|[来源1](https://www.dunloptennis.net/storage/catalog2025.pdf) / [来源2](https://us.dunlopsports.com/tennis-strings.html) 日本官方2025目录明确毫米线径；为目录覆盖，不保证中国所有渠道当前有货。|
|Tecnifibre|RAZOR SPIN|1.15, 1.20, 1.25, 1.30|[来源1](https://www.tecnifibre.jp/products/string/razor-spin/) / [来源2](https://www.tecnifibre.jp/wp/wp-content/uploads/04d44dcc1be25c5336fe15c9144338ab.pdf) 2026-05-16日本上市，官方新品公告明确毫米线径；材质含PU添加，不表示纯聚酯。|
|Tecnifibre|NRG2|1.24|[来源1](https://www.tecnifibre.jp/products/string/nrg2/) / [来源2](https://www.tennis-warehouse.com/Tecnifibre_NRG2_17_124_String_Black/descpageACTFUSA-NRG17BK.html) 本批仅确认1.24；官网产品编码中的132不直接转换成毫米规格。含PU。|
|Tecnifibre|TGV|1.25, 1.30|[来源1](https://www.tecnifibre.jp/products/string/tgv/) / [来源2](https://www.tecnifibre.com/collections/cordages-de-tennis-multifilament/products/bobine-200m-tgv-black) 毫米线径由品牌全球商店索引确认，材料含PU。1.35尚未纳入。|
|Tecnifibre|XR3|1.25, 1.30|[来源1](https://www.tecnifibre.com/en/products/garniture-xr3-naturel) 含三个单丝芯和多股结构；主材质归一为聚酰胺，含PU与硅表面处理。|
|Tecnifibre|MULTIFEEL|1.25, 1.30|[来源1](https://www.tecnifibre.com/en-ue/products/bobine-multifeel-black?variant=54757820236120) 聚酰胺中心单丝加多股结构，含PU；不将中心单丝误分类为整体单丝。|
|Tecnifibre|TRIAX|1.28|[来源1](https://www.tecnifibre.com/en/products/garniture-triax-naturel) / [来源2](https://www.tecnifibre.jp/products/string/triax/) 日本官网明确尼龙、聚酯、PU混合；全球页明确1.28。1.33尚未纳入，不能按产品编号推断。|
|Tecnifibre|DURAMIX|1.25|[来源1](https://www.tecnifibre.com/en-uee/products/garniture-duramix-black) / [来源2](https://www.tecnifibre.jp/products/string/duramix/) 材料含尼龙、聚酯、PU。全球标题写monofilament与日本结构描述冲突，采用日本规格明确的多股结构；仅收毫米明示的1.25。|
|Luxilon|Element IR Soft|1.27|[来源1](https://www.wilson.com/de-de/products/element-ir-soft-127-tennis-string-set-wr8309200?variant=53965119029562) / [来源2](https://catalogs.asbsports.com/vendors/pdfs/Wilson-Racquet_Sports-SS_2026.pdf) 官方商品页确认1.27，材质由品牌2026目录Polyester分类交叉确认；与Element分开建型号。|

## 数据库写入与核验

- 新增按 (brand, model) 唯一键定位；规格按 (string_id, diameter_mm) 去重，不指定数据库ID，不删除已有型号和规格。重复执行不会新增重复记录。
- 测试环境新增项 published=true，可在小程序搜索。资料JSON published=false是采集阶段标记；2026-10-27生产与测试库均已启用，不以该历史标记推断实时发布状态。
- X-ONE BIPHASE仅在现有material为空时补为polyamide并更新来源，不改个人穿线快照。
- 暂未知材质共6款：HEAD Lynx、Hawk、Velocity MLT、Lynx Touch、Hawk Power，以及Wilson NXT。字段留空不影响按名称与线径记录；不对用户展示推测材料。
- DB只存一个source_url；完整多来源、核对日期、结构、缺项和备注保留在版本化资料中，未新增数据库字段。

## 生产交接

1. 生产部署必须包含球线库及穿线结构迁移；确认API与小程序版本一致。
2. 合并第一批items、第二批items和existingItemUpdates，复核6条未知材料及需要覆盖的地区SKU。不要直接复制测试库ID。
3. 使用(brand, model)和(string_id, diameter_mm)在事务中做增量导入；新条目发布前核对来源、名称、毫米规格。数据源checked不代表已取得商品图片或文案授权，本库不提供这些内容。
4. 验收品牌筛选、型号/别名搜索、线径选择、手动输入、混线保存、历史快照回看。生产导入与发布需在明确指定生产环境后另行执行。

本次仅扩充基础目录，API及数据库结构不变。


## 2026-10-27 生产同步结果

- 用户明确要求同步线上，沿用不备份要求。导入前线上string_library、string_variants均为空。
- 在事务内按品牌/型号和型号ID/毫米线径增量插入，71款、184规格全部published=true；不复制测试库ID，不修改用户穿线记录。
- X-ONE BIPHASE按第二批核实结果写入polyamide；其余6款未知材料留空。
- 提交前及提交后逐项比对名称、中文名、别名、来源URL、材料、发布状态与全部毫米规格，均通过。
- 本次仅同步目录数据；没有部署后端或小程序代码，也未做线上客户端端到端验收。
