# 首批主流网球线资料（40款）

整理日期：2026-10-21。用于穿线选择库，属于编辑选择的常见型号覆盖，未经销量统计，不代表排名或适配推荐。

共 40 个型号、102 个已核对毫米线径规格，覆盖 7 个品牌。基础字段尚有缺项的 4 款；逐条见备注。

数据文件：[mainstream.v1.json](../data/tennis-strings/mainstream.v1.json)。首批40款已导入本地测试库并启用；JSON保留采集阶段 published=false 标记；2026-10-27已将两批资料共71款、184规格同步生产并启用。2026-10-27新增第二批后测试库为71款、184规格，见[增量资料与上线交接](tennis-strings-expansion.md)。

## 品牌与型号

| 品牌 | 型号数 |
|---|---:|
| Yonex | 10 |
| Solinco | 6 |
| Tecnifibre | 6 |
| Luxilon | 5 |
| HEAD | 5 |
| Babolat | 6 |
| Wilson | 2 |

## 逐项清单

材质字段与结构字段分开；“待核实”不是材质分类。线径均为 mm，表示已查到的规格子集，不表示全部可售规格。

| 品牌 | 型号 | 材质 | 结构 | 本次核对线径 | 来源与说明 |
|---|---|---|---|---|---|
| Yonex | POLYTOUR PRO | 聚酯 | 单丝 | 1.15, 1.20, 1.25, 1.30 | [来源1](https://www.yonex.com/tennis-strings) |
| Yonex | POLYTOUR REV | 聚酯 | 单丝 | 1.20, 1.25, 1.30 | [来源1](https://www.yonex.com/tennis-strings) |
| Yonex | POLYTOUR FIRE | 聚酯 | 单丝 | 1.20, 1.25 | [来源1](https://www.yonex.com/tennis-strings)；官网标题包含130，但Gauge正文只列120/125。本批仅收录正文明确规格，130待进一步核对。 |
| Yonex | POLYTOUR SPIN | 聚酯 | 单丝 | 1.20, 1.25 | [来源1](https://www.yonex.com/tennis-strings) |
| Yonex | POLYTOUR AIR | 聚酯 | 单丝 | 1.25 | [来源1](https://www.yonex.com/tennis-strings)；官网另注明HR-Elastomer；material为主材质归一值，不表示纯聚酯。 |
| Yonex | POLYTOUR DRIVE | 聚酯 | 单丝 | 1.25 | [来源1](https://www.yonex.com/tennis-strings) |
| Yonex | REXIS FEEL | 聚酰胺/尼龙 | 多股 | 1.25, 1.30 | [来源1](https://www.yonex.com/tennis-strings)；官网另列多种聚合物添加组分，仅归一主材质。 |
| Yonex | REXIS COMFORT | 聚酰胺/尼龙 | 多股 | 1.25, 1.30 | [来源1](https://www.yonex.com/tennis-strings)；官网为Nylon并含FORTIMO。 |
| Yonex | REXIS SPEED | 聚酰胺＋聚酯 | 多股 | 1.25, 1.30 | [来源1](https://www.yonex.com/tennis-strings)；混合材质，不能归为纯尼龙或纯聚酯。 |
| Yonex | DYNAWIRE | 聚酰胺/尼龙 | 单丝 | 1.25, 1.30 | [来源1](https://www.yonex.com/tennis-strings)；官网为Nylon monofilament并含Metal Film，按其原始结构分类记录。 |
| Solinco | Hyper-G | 共聚酯 | 单丝 | 1.05, 1.10, 1.15, 1.20, 1.25, 1.30, 1.35 | [来源1](https://solincosports.com/equipment/hyper-g/) |
| Solinco | Hyper-G Soft | 共聚酯 | 单丝 | 1.15, 1.20, 1.25, 1.30 | [来源1](https://solincosports.com/equipment/hyper-g-soft/)；Soft、Round与普通版分别建型号，不能合并。 |
| Solinco | Hyper-G Round | 共聚酯 | 单丝 | 1.15, 1.20, 1.25, 1.30 | [来源1](https://solincosports.com/equipment/hyper-g-round/)；Soft、Round与普通版分别建型号，不能合并。 |
| Solinco | Tour Bite | 共聚酯 | 单丝 | 1.05, 1.10, 1.15, 1.20, 1.25, 1.30, 1.35 | [来源1](https://solincosports.com/equipment/tour-bite/) |
| Solinco | Tour Bite Soft | 共聚酯 | 单丝 | 1.15, 1.20, 1.25, 1.30 | [来源1](https://solincosports.com/equipment/tour-bite-soft/)；Soft、Round与普通版分别建型号，不能合并。 |
| Solinco | Confidential | 共聚酯 | 单丝 | 1.15, 1.20, 1.25, 1.30 | [来源1](https://solincosports.com/equipment/confidential/) |
| Tecnifibre | RAZOR SOFT | 聚酯 | 单丝 | 1.20, 1.25, 1.30 | [来源1](https://www.mytecnifibre.com/site/img/Home/Tecnifibre%20USA%20SS26%20-%20Tennis%2C%20Padel%2C%20Squash%2C%20and%20Pickleball%20Catalogue.pdf)；来源为厂商品牌美国2026目录18–19页；颜色可售组合未建模，发布前结合原表复核。 |
| Tecnifibre | RAZOR CODE | 聚酯 | 单丝 | 1.20, 1.25, 1.30 | [来源1](https://www.mytecnifibre.com/site/img/Home/Tecnifibre%20USA%20SS26%20-%20Tennis%2C%20Padel%2C%20Squash%2C%20and%20Pickleball%20Catalogue.pdf)；来源为厂商品牌美国2026目录18–19页；颜色可售组合未建模，发布前结合原表复核。 |
| Tecnifibre | BLACK CODE | 共聚酯 | 单丝 | 1.18, 1.24, 1.28, 1.32 | [来源1](https://www.mytecnifibre.com/site/img/Home/Tecnifibre%20USA%20SS26%20-%20Tennis%2C%20Padel%2C%20Squash%2C%20and%20Pickleball%20Catalogue.pdf)；来源为厂商品牌美国2026目录18–19页；颜色可售组合未建模，发布前结合原表复核。 |
| Tecnifibre | RED CODE | 共聚酯 | 单丝 | 1.20, 1.25, 1.30 | [来源1](https://www.mytecnifibre.com/site/img/Home/Tecnifibre%20USA%20SS26%20-%20Tennis%2C%20Padel%2C%20Squash%2C%20and%20Pickleball%20Catalogue.pdf)；来源为厂商品牌美国2026目录18–19页；颜色可售组合未建模，发布前结合原表复核。 |
| Tecnifibre | ICE CODE | 共聚酯 | 单丝 | 1.20, 1.25, 1.30 | [来源1](https://www.tecnifibre.com/en/products/garniture-ice-code-blanc) |
| Tecnifibre | X-ONE BIPHASE | 待核实 | 多股 | 1.24, 1.30 | [来源1](https://www.tecnifibre.com/en/collections/cordages-tennis/products/garniture-x-one-biphase-black)；所查产品页确认多股结构及线径，没有充分确认化学材质；保留null，不能把multifilament填入material。 |
| Luxilon | ALU Power | 聚酯 | 单丝 | 1.25 | [来源1](https://www.wilson.com/en-us/products/alu-power-125-tennis-string-set-rz99203)；仅收录本次产品页确认的1.25，并非该型号完整线径目录；官网还描述其专有聚合物与铝组分。 |
| Luxilon | ALU Power Rough | 聚酯 | 单丝 | 1.25 | [来源1](https://id.wilson.com/en/products/bb-alupower-rough-125-str)；本次仅确认125，不将介绍中提到的130自动视为已确认SKU。 |
| Luxilon | ALU Power Soft | 聚酯 | 单丝 | 1.25 | [来源1](https://www.wilson.com/en-us/products/alu-power-soft-125-tennis-string-set-rz99102) |
| Luxilon | 4G | 聚酯 | 单丝 | 1.25 | [来源1](https://www.wilson.com/en-us/products/4g-125-tennis-string-set-rz99710)；本次只确认125；4G Soft不能合并进4G。 |
| Luxilon | Element | 聚酯 | 待补 | 1.25 | [来源1](https://www.wilson.com/en-us/collections/tennis?irgwc=1) / [来源2](https://www.wilson.com/en-us/collections/luxilon-tennis-string)；型号125由官方在售目录确认，聚酯分类由官方品牌集合页确认；具体结构和其它线径待补。 |
| HEAD | Lynx Tour | 共聚酯 | 单丝 | 1.25, 1.30 | [来源1](https://www.head.com/en_GB/product/lynx-tour-4) |
| HEAD | Lynx | 待核实 | 单丝 | 1.20, 1.25, 1.30 | [来源1](https://www.head.com/en_US/product/lynx-6)；当前所查规格页明确单丝和线径，化学材质未在本次证据中明确，待补。 |
| HEAD | Hawk | 待核实 | 单丝 | 1.20, 1.25, 1.30 | [来源1](https://www.head.com/en_CA/product/hawk-set-281103)；单丝和线径已确认；化学材质待补。 |
| HEAD | Hawk Touch | 共聚酯 | 单丝 | 1.15, 1.20, 1.25, 1.30 | [来源1](https://www.head.com/fr_CA/product/hawk-touch-3) / [来源2](https://www.tennis-warehouse.com/Head_Hawk_Touch_17_125_String_Anthracite/descpageACHEAD-HAWKT17.html)；材质由TW具体商品规格交叉核对。 |
| HEAD | Hawk Tour rPET | 共聚酯 | 单丝 | 1.25, 1.30 | [来源1](https://www.head.com/fr_BE/product/hawk-tour-rpet-281125)；厂商规格为再生PET共聚酯。 |
| Babolat | RPM Blast | 共聚酯 | 单丝 | 1.20, 1.25, 1.30, 1.35 | [来源1](https://www.babolat.com/us/rpm-blast-12m/241101.html) / [来源2](https://www.tennis-warehouse.com/RPM_Blast_/catpage-RPMFAM.html)；线径由TW该系列目录交叉确认；新出的蓝色1.27暂不并入普通型号，待颜色/配方版本核对。 |
| Babolat | RPM Team | 共聚酯 | 单丝 | 1.25, 1.30 | [来源1](https://www.babolat.com/gb/rpm-team-12m/241108.html) / [来源2](https://www.tennis-warehouse.com/descpage.html?img=6&pcode=P3LY05)；官网Gauge选项未展开；线径由TW商品索引交叉确认。 |
| Babolat | RPM Rough | 共聚酯 | 单丝 | 1.25 | [来源1](https://www.babolat.com/us/rpm-rough-200m/201-243140.html) / [来源2](https://www.tennis-warehouse.com/Babolat_RPM_Rough_17_125_String/descpageACBAB-BRPMBRS.html?color=GR)；本次仅确认TW的1.25商品规格，非完整线径范围。 |
| Babolat | RPM Soft | 聚酰胺/尼龙 | 待补 | 1.25, 1.30 | [来源1](https://www.babolat.com/us/tennis/strings.html) / [来源2](https://www.tennis-warehouse.com/descpage.html?img=6&pcode=P3LY05)；重要：RPM Soft材质为聚酰胺，不能因RPM前缀归成聚酯。线径由TW商品索引交叉确认。 |
| Babolat | Xcel | 聚酰胺/尼龙 | 多股 | 1.25 | [来源1](https://www.babolat.com/us/tennis/strings.html) / [来源2](https://www.tennis-warehouse.com/Babolat_Xcel_17125_String/descpageACBAB-BXCL17.html)；本次仅确认1.25，非全部规格。 |
| Babolat | Touch VS | 天然肠线 | 天然肠线 | 1.30 | [来源1](https://www.babolat.com/us/rpm-blast-125-touch-vs-130/281037.html?dwvar_281037_COLOR_DESCRIPTION_ERP=105)；以官网组合装确认Touch VS为天然肠线且130；拆出单根线材型号，组合装本身不单独建库。 |
| Wilson | NXT | 待核实 | 多股 | 1.30 | [来源1](https://www.tennis-warehouse.com/Wilson_NXT_16_130_String_Black/descpageACWILSON-NXT16B.html) / [来源2](https://www.wilson.com/en-us/collections/tennis)；该条规格与结构来源为零售商TW，官网目录仅确认产品名；化学材质不作推断。 |
| Wilson | Sensation | 聚酰胺/尼龙 | 多股 | 1.25, 1.30 | [来源1](https://www.wilson.com/en-gb/products/sensation-16-tennis-string-200m-reel-rz91100)；官方详情同时列16/17和1.30/1.25，按型号汇总两个线径。 |

## 数据处理约定

- 不收价格、推荐磅数、舒适度分数或推算性能；不复制商品图、营销文案或评论。
- 英文型号作标准名，中文型号 nameZh 暂空。搜索别名为编辑补充，不以昵称建立重复商品。
- 大盘/小包装合并；Soft、Rough、Round、Tour等不同产品分别保留。
- 厂商/地区页与零售商可有SKU差异。本批只确认已列规格，未确认的不推定不存在。
- 每项附核对日期、来源、缺项及备注。来源checked不等于已取得图片、文案或数据集再发布授权。
- 首次导入前需要复核missingFields、来源质量与使用范围；本文件不能直接作为现有Go DTO请求发送。JSON aliases为数组，落库需由导入工具转换为当前aliases文本，材质null也需按实现处理；不要伪造variantId。

## 优先补核事项

1. HEAD Lynx、Hawk、Tecnifibre X-ONE BIPHASE、Wilson NXT的化学材质证据。
2. Babolat Xcel除已确认1.25之外的线径明细。
3. Yonex POLYTOUR FIRE：标题130与正文规格不一致，先不补130。
4. Luxilon Element及Babolat RPM Soft的具体结构；当前可保存未知结构。
5. Luxilon ALU系列其它线径、RPM Blast蓝色1.27的版本/颜色关系；不得把推测作为规格。
6. 中文别名可根据真实搜索反馈补充；PTP只指POLYTOUR PRO，不扩散到其它POLYTOUR型号。
