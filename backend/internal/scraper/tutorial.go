// backend/internal/scraper/tutorial.go
package scraper

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// SeedTutorialSite 生成一起出趣教程写入站点
func SeedTutorialSite(db *gorm.DB, siteID uint) error {
	if err := ClearPages(db, siteID); err != nil {
		return fmt.Errorf("clear pages: %w", err)
	}

	type pageSpec struct {
		parentSlug string
		slug       string
		title      string
		md         string
	}
	pages := []pageSpec{
		{"", "intro", "简介", tutorialIntro},
		{"intro", "what", "一起出趣是什么", tutorialWhat},
		{"", "quickstart", "新手使用教程", tutorialQuickstart},
		{"", "features", "核心功能", ""},
		{"features", "plan", "共享行程计划", tutorialFeaturePlan},
		{"features", "checklist", "共享出趣清单", tutorialFeatureChecklist},
		{"features", "location", "实时好友定位", tutorialFeatureLocation},
		{"features", "memory", "共享回忆", tutorialFeatureMemory},
		{"features", "news", "自动更新动态消息", tutorialFeatureNews},
		{"features", "bill", "共享出游账单", tutorialFeatureBill},
		{"", "highlights", "软件亮点", tutorialHighlights},
		{"", "faq", "常见问题", tutorialFAQ},
	}

	// 先建一级，再建二级
	slugToID := map[string]uint{}
	sortIdx := 0
	for _, p := range pages {
		if p.parentSlug == "" {
			id, err := InsertPage(db, siteID, nil, p.slug, p.title, p.md, sortIdx)
			if err != nil {
				log.Printf("insert %q: %v", p.title, err)
				continue
			}
			slugToID[p.slug] = id
			sortIdx++
		}
	}
	childSort := 0
	for _, p := range pages {
		if p.parentSlug != "" {
			parentID, ok := slugToID[p.parentSlug]
			if !ok {
				continue
			}
			pid := parentID
			_, err := InsertPage(db, siteID, &pid, p.slug, p.title, p.md, childSort)
			if err != nil {
				log.Printf("insert child %q: %v", p.title, err)
				continue
			}
			childSort++
		}
	}
	return nil
}

const tutorialIntro = `# 一起出趣 简介

**一起出趣** 是一款专为旅行爱好者打造的社交旅行应用，致力于搭建一个便捷的旅行规划与分享平台。

用户能够在平台上发布个人旅行计划，邀约志同道合的朋友共同出游，也可以浏览他人发布的行程，挑选合适的同行伙伴。同时支持实时记录旅行中的精彩时刻，分享照片、视频以及游记，与世界各地的旅行者展开互动。

## 适用场景

- 同城社交出行
- 户外活动结伴
- 旅行计划多人协作
- 旅途回忆共享

> 本教程内容由公开资料整理生成，可在后台编辑修改。
`

const tutorialWhat = `# 一起出趣是什么

一起出趣是一款集旅行计划、行程分享、交友于一体的综合性旅游平台，让每一次旅行都充满无限可能。

## 主要功能概览

- **旅伴**：寻找志同道合的旅伴
- **社区交流**：查看附近旅行者动态
- **发布旅行照片**：记录旅途见闻
- **发布指南**：分享个性化旅行指南
- **发布点评**：景点点评与商家推荐
- **特色体验**：发红包获得实时帮助、网红旅游目的地签到和记录

## 筛选与收藏

- 按关键字和排序搜索你的理想活动
- 随时随地完成预订或添加到愿望清单
`

const tutorialQuickstart = `# 新手使用教程

5 步快速上手一起出趣小程序/app。

## 第 1 步：进入"旅行"板块

点开"旅行"板块，即可着手创建个人专属的旅行计划。

## 第 2 步：进入"密友"板块

点开"密友"板块，添加感兴趣的趣友并与之建立联系。

## 第 3 步：创建新旅行

于"旅行"界面，点击"添加新旅行"以创建一个全新的旅行项目。

## 第 4 步：通过暗号加入他人旅行

若想加入他人的旅行，点击"暗号"即可轻松加入。

## 第 5 步：复制口令添加趣友

把好友提供的口令复制下来，直接添加进自己的趣友列表。

> 提示：暗号和口令是结伴出游的核心机制，务必保管好你分享的口令。
`

const tutorialFeaturePlan = `# 共享行程计划

我们提供便捷的多人协作功能，让您和伙伴们轻松创建并编辑行程计划。

## 使用方法

1. 在旅行详情页点击"编辑行程"
2. 添加每日活动安排
3. 邀请伙伴共同编辑
4. 地图模式一键呈现，让行程规划更直观

## 亮点

- 多人实时协作，所见即所得
- 地图视图直观展示每日路线
- 行程变动自动同步给所有成员
`

const tutorialFeatureChecklist = `# 共享出趣清单

出游必备事项不再遗漏！

## 功能说明

与伙伴们共同创建、编辑并跟踪待办事项，谁负责钥匙、谁忘了证件，一目了然，确保旅途无忧。

## 推荐清单模板

- 身份证 / 护照
- 钥匙
- 充电器 / 移动电源
- 常用药品
- 现金 / 银行卡
- 帐篷（如需露营）
- 车辆安排确认
`

const tutorialFeatureLocation = `# 实时好友定位

想要共享信息？随时开启，即刻掌握伙伴们的实时位置。

## 使用场景

- 轻松查看谁已到达集合点，谁还在路上
- 告别无谓的等待
- 途中安全提醒

## 隐私说明

- 位置共享可随时关闭
- 仅在同行人之间可见
- 不记录历史轨迹
`

const tutorialFeatureMemory = `# 共享回忆

出游照片不再孤单，直接关联活动并共享给伙伴们。共同记录美好瞬间，珍藏快乐回忆。

## 使用方法

1. 在活动详情页点击"上传照片"
2. 选择关联的活动
3. 自动共享给所有同行伙伴
4. 伙伴可点赞、评论
`

const tutorialFeatureNews = `# 自动更新动态消息

行程计划或待办事项有任何变动？无需逐一通知，动态消息将自动推送至所有相关成员，确保信息同步，减少沟通成本。

## 消息类型

- 行程变更
- 待办完成
- 新成员加入
- 照片上传
- 账单更新
`

const tutorialFeatureBill = `# 共享出游账单

旅途中，实现多人分账、支出平摊并自动结算，实时查看成员结算方案以及剩余应付款项。

## 使用流程

1. 在旅行详情页打开"账单"
2. 录入每笔支出（金额、付款人、参与人）
3. 系统自动计算每人应付
4. 行程结束一键生成分账方案
`

const tutorialHighlights = `# 软件亮点

## 探索

让你脑洞大开，探索旅行目的地的灵感指南。

## 分类

轻松找到符合你兴趣与预算的各类项目。

## 应用场景

同城社交出行、户外活动。

## 主要功能

旅伴、社区交流、查看附近旅行的人、发布旅行照片、发布指南、发布点评、商家推荐。

## 特色体验

可以发红包获得实时帮助，可以在网红旅游目的地签到和记录。

## 筛选

按关键字和排序搜索你的理想活动。

## 收藏

随时随地完成预订或添加到愿望清单。
`

const tutorialFAQ = `# 常见问题

## Q: 一起出趣是免费的吗？

A: 应用本身免费下载使用，部分高级功能可能需要应用内购买。

## Q: 支持哪些系统？

A: 当前主要支持 Android 系统（4.3 以上）。

## Q: 如何找到合适的旅伴？

A: 可以通过"密友"板块添加感兴趣的趣友，也可以浏览他人发布的行程，挑选合适的同行伙伴。

## Q: 暗号和口令有什么区别？

A: 暗号用于加入他人的旅行项目，口令用于直接添加对方到自己的趣友列表。

## Q: 位置共享会泄露隐私吗？

A: 位置共享可随时关闭，仅在同行人之间可见，不记录历史轨迹。

## Q: 账单分账支持几人？

A: 支持旅行中所有成员参与分账，系统自动计算每人应付金额。

> 本 FAQ 为通用内容整理，具体以应用内实际说明为准。
`
