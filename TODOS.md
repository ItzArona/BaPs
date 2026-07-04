# TODOS — 未实现的 Excel 数据表

本文档由 workflow 日志自动生成,记录 `generate_excel.go` 在构建
`data/Excel.bin` 时因 `pb.Excel` proto 无对应字段而跳过的所有数据表。

- **数据来源**: CI run `28692957858` (upload-data job) 的 `Skipping XxxExcel` 日志
- **已实现**: 62 张表(已加载进 `pb.Excel`)
- **未实现**: 411 张表(按下文分类列出,共 411 项)
- **生成时间**: 2026-07-04

> 注意: 下表中的 "未实现" 仅表示后端 `pb.Excel` 结构体里没有对应字段,
> `generate_excel.go` 跳过了该 JSON。很多表是纯客户端数据(动画/音效/UI),
> 服务端永远不需要;也有不少是后端应该支持但尚未适配的玩法表。每个分类
> 里已用备注粗略区分,实际是否需要适配请结合玩法逻辑判断。

---

## 1. 客户端资源 / 动画 / 音效 / UI（服务端通常不需要）

这些表只服务于客户端渲染,服务端不读取,一般可永久忽略。

```
AnimationBlendTable
AnimatorDataTable
AudioAnimatorExcel
BGMExcel
BGMRaidExcel
BGMUIExcel
CameraExcel
CharacterIllustCoordinateExcel
CharacterVictoryInteractionExcel
CombatEmojiExcel
HpBarAbbreviationExcel
LoadingImageExcel
MiniGameAudioAnimatorExcel
SoundUIExcel
SpineLipsyncExcel
VideoExcel
VoiceCommonExcel
VoiceExcel
VoiceLogicEffectExcel
VoiceRoomExceptionExcel
VoiceSpineExcel
VoiceTimelineExcel
VoiceRoomExceptionExcel
```

## 2. 本地化文本表（服务端按需引用,通常不整表加载）

服务端一般只通过 `LocalizeEtcId` 间接引用文案,这些整张本地化表
通常不需要进 `pb.Excel`。

```
LocalizeCCGExcelTable
LocalizeCharProfileChangeExcel
LocalizeCharProfileExcel
LocalizeErrorExcel
LocalizeEtcExcel
LocalizeEtcExcelTable
LocalizeExcel
LocalizeFieldExcelTable
LocalizeGachaShopExcel
LocalizeSNSExcelTable
LocalizeSkillExcel
```

## 3. 场景 / 剧情 / 演出（服务端按需引用）

剧情演出相关的脚本、表情、立绘坐标、转场等,多为客户端演出数据。

```
ScenarioBGEffectExcel
ScenarioBGNameExcel
ScenarioCharacterEmotionExcel
ScenarioCharacterNameExcel
ScenarioCharacterSituationSetExcel
ScenarioContentCollectionExcel
ScenarioEffectExcel
ScenarioModeSpoilerPopupExcel
ScenarioReplayExcelTable
ScenarioResourceInfoExcel
ScenarioScriptExcel
ScenarioScriptFieldExcelTable
ScenarioTransitionExcel
CharacterDialogBattlePassExcel
CharacterDialogEmojiExcel
CharacterDialogEmojiExcelTable
CharacterDialogEventExcel
CharacterDialogExcel
CharacterDialogFieldExcelTable
CharacterDialogSubtitleExcel
CharacterVoiceExcel
CharacterVoiceSubtitleExcel
TutorialCharacterDialogExcel
MomotalkScheduleSpoilerPopupExcel
ContentsScenarioExcel
```

## 4. 活动 / EventContent（玩法核心,大部分应实现）

活动玩法相关表,后端如要支持活动循环、商店、卡牌、宝藏等需要适配。

### 4.1 活动通用 / 季节 / 通知
```
EventContentArchiveBannerOffsetExcel
EventContentChangeExcel
EventContentChangeScenarioExcel
EventContentLobbyMenuExcel
EventContentMissionExcel
EventContentNotifyExcel
EventContentPlayGuideExcel
EventContentScenarioExcel
EventContentSeasonExcel
EventContentShopInfoExcel
EventContentSpecialOperationsExcel
EventContentSpineDialogOffsetExcel
EventContentSpineDisplayPeriodExcel
EventContentSpoilerPopupExcel
EventContentZoneExcel
EventContentCurrencyItemExcel
EventContentMiniEventShortCutExcel
EventContentMiniEventTokenExcel
ConstEventCommonExcelTable
```

### 4.2 活动关卡 / 奖励
```
EventContentStageExcel
EventContentStageRewardExcel
EventContentStageTotalRewardExcel
EventContentLocationExcel
EventContentLocationRewardExcel
```

### 4.3 活动商店
```
EventContentShopExcel
EventContentCardExcel
EventContentCardShopExcel
EventContentCardShopModifyExcel
EventContentBoxGachaElementExcelTable
EventContentBoxGachaManageExcel
EventContentBoxGachaShopExcel
EventContentFortuneGachaExcel
EventContentFortuneGachaModifyExcel
EventContentFortuneGachaShopExcel
```

### 4.4 活动子玩法:线索搜索 / 集中 / 宝藏 / 骰子赛 / 见面会
```
EventContentClueExcel
EventContentClueSearchExcel
EventContentClueSearchRewardExcel
EventContentClueSearchRoundExcel
EventContentCollectionExcel
EventContentConcentrationCardExcel
EventContentConcentrationExcel
EventContentConcentrationRewardExcel
EventContentConcentrationVoiceExcel
EventContentTreasureExcel
EventContentTreasureRewardExcel
EventContentTreasureRoundExcel
EventContentTreasureCellRewardExcel
EventContentDiceRaceEffectExcel
EventContentDiceRaceExcel
EventContentDiceRaceNodeExcel
EventContentDiceRaceProbExcel
EventContentDiceRaceTotalRewardExcel
EventContentMeetupExcel
EventContentMeetupInfoExcel
EventContentBuffExcel
EventContentBuffGroupExcel
EventContentCharacterBonusExcel
EventContentDebuffRewardExcel
```

## 5. 小游戏 / Minigame（独立玩法模块）

### 5.1 CCG 卡牌
```
MinigameCCGCardExcel
MinigameCCGCharacterExcel
MinigameCCGEnemyExcel
MinigameCCGEnemyGroupExcel
MinigameCCGInfoExcel
MinigameCCGLevelExcel
MinigameCCGLevelNodeExcel
MinigameCCGLevelStageExcel
MinigameCCGLogicEffectExcel
MinigameCCGOpenDialogExcel
MinigameCCGPerkExcel
MinigameCCGRewardCardExcel
MinigameCCGRewardCardRateExcel
MinigameCCGRewardItemExcel
MinigameCCGSkillExcel
MinigameCCGStartDeckCardExcel
MinigameCCGStartDeckCharacterExcel
ConstMinigameCCGExcelTable
```

### 5.2 TBG 桌游
```
MinigameTBGDiceExcel
MinigameTBGEncounterExcel
MinigameTBGEncounterOptionExcel
MinigameTBGEncounterRewardExcel
MinigameTBGItemExcel
MinigameTBGObjectExcel
MinigameTBGSeasonExcel
MinigameTBGThemaExcel
MinigameTBGThemaRewardExcel
MinigameTBGVoiceExcel
ConstMinigameTBGExcelTable
```

### 5.3 Dream 收集
```
MiniGameDreamCollectionScenarioExcel
MiniGameDreamCollectionScenarioExcelTable
MiniGameDreamDailyPointExcel
MiniGameDreamEndingExcel
MiniGameDreamEndingRewardExcel
MiniGameDreamInfoExcel
MiniGameDreamParameterExcel
MiniGameDreamReplayScenarioExcel
MiniGameDreamScheduleExcel
MiniGameDreamScheduleResultExcel
MiniGameDreamTimelineExcel
MinigameDreamVoiceExcel
```

### 5.4 RoadPuzzle 接龙
```
MiniGameRoadPuzzleInfoExcel
MiniGameRoadPuzzleRailSetRewardExcel
MiniGameRoadPuzzleRewardExcel
MiniGameRoadPuzzleVoiceExcel
MinigameRoadPuzzleAdditionalRewardExcel
MinigameRoadPuzzleMapExcel
MinigameRoadPuzzleMapTileExcel
MinigameRoadPuzzleRailTileExcel
MinigameRoadPuzzleRoadRoundExcel
ConstMinigameRoadPuzzleExcelTable
```

### 5.5 Shooting 射击
```
MiniGameShootingCharacterExcel
MiniGameShootingGeasExcel
MiniGameShootingStageExcel
MiniGameShootingStageRewardExcel
ConstMiniGameShootingExcelTable
```

### 5.6 Defense 塔防
```
MiniGameDefenseCharacterBanExcel
MiniGameDefenseFixedStatExcel
MiniGameDefenseFixedStatExcelTable
MiniGameDefenseInfoExcel
MiniGameDefenseStageExcel
```

### 5.7 Rhythm 节奏 / 其他
```
MiniGameRhythmBgmExcel
MiniGameRhythmExcel
MinigameCardExcelTable
MiniGameMissionExcel
MiniGamePlayGuideExcel
```

## 6. 征服战 / Conquest（大型活动玩法）

```
ConquestCalculateExcel
ConquestCameraSettingExcel
ConquestErosionExcel
ConquestErosionUnitExcel
ConquestEventExcel
ConquestGroupBonusExcel
ConquestGroupBuffExcel
ConquestMapExcel
ConquestObjectExcel
ConquestPlayGuideExcel
ConquestProgressResourceExcel
ConquestRewardExcel
ConquestStepExcelTable
ConquestTileExcel
ConquestUnexpectedEventExcel
ConquestUnitExcel
ConstConquestExcelTable
```

## 7. 世界Raid / WorldRaid & 互动世界Raid

```
WorldRaidBossGroupExcel
WorldRaidConditionExcel
WorldRaidFavorBuffExcel
WorldRaidSeasonManageExcel
WorldRaidStageExcel
WorldRaidStageRewardExcel
InteractiveWorldRaidArcadeMachineExcel
InteractiveWorldRaidBossGroupExcel
InteractiveWorldRaidCarrierExcel
InteractiveWorldRaidCarrierExcelTable
InteractiveWorldRaidCarrierMapExcel
InteractiveWorldRaidCarrierRecipeExcel
InteractiveWorldRaidConditionExcel
InteractiveWorldRaidSeasonManageExcel
InteractiveWorldRaidSkillDescriptionListExcel
InteractiveWorldRaidStageExcel
InteractiveWorldRaidStatusPresetExcel
```

## 8. 角色相关扩展表（部分可能需要适配）

```
CharacterAIExcel
CharacterAcademyTagsExcel
CharacterCalculationLimitExcel
CharacterCombatSkinExcel
CharacterExcel                  # 注意: 已实现的是 CharacterExcelTable(Excel.zip), 此为 ExcelDB 版
CharacterGearLevelExcel
CharacterLevelExcel             # 注意: 已实现的是 CharacterLevelExcelTable
CharacterLevelStatFactorExcel
CharacterPotentialRewardExcel
CharacterSkillListExcel         # 注意: 已实现的是 CharacterSkillListExcelTable
CharacterStatExcel              # 注意: 已实现的是 CharacterStatExcelTable
CharacterStatLimitExcel
CharacterStatsDetailExcel
CharacterStatsTransExcel
CharacterTranscendenceExcel
CharacterWeaponExpBonusExcel
CharacterWeaponLevelExcel
CostumeExcel
PersonalityExcel
TerrainAdaptationFactorExcel
BulletArmorDamageFactorExcel
BattleLevelFactorExcel
StatLevelInterpolationExcel
GrowthScoreCalculationExcel
```

## 9. 战斗 / 战术模拟器 / Tactic

```
TacticArenaSimulatorSettingExcelTable
TacticDamageSimulatorSettingExcelTable
TacticEntityEffectFilterExcel
TacticSimulatorSettingExcelTable
TacticSkipExcel
TacticTimeAttackSimulatorConfigExcelTable
TacticalSupportSystemExcel
BossExternalBTExcel
BossPhaseExcelTable
FixedStrategyExcel
StrategyObjectBuffDefineExcel
CampaignStrategyObjectExcel
InformationStrategyObjectExcel
StoryStrategyExcel
ClearDeckRuleExcelTable
ObstacleFireLineCheckExcel
ObstacleStatExcel
FloaterCommonExcel
GroundExcel
GroundModuleRewardExcel
FormationLocationExcel
EchelonConstraintExcel
FixedEchelonSettingExcel
ConstCombatExcelTable
ConstStrategyExcelTable
UnderCoverStageExcel
```

## 10. 商店扩展 / Shop & Product

```
ShopCashExcel
ShopCashScenarioResourceInfoExcel
ShopFilterClassifiedExcel
ShopFreeRecruitExcel
ShopFreeRecruitPeriodExcel
ShopRecruitDirectingExcel
ShopRecruitExcel
ShopRecruitSettingExcel
ShopTabGroupExcel
ProductAutoSelectionGroupExcel
ProductBattlePassExcel
ProductBattlePassExcelTable
ProductDailyRecordExcel
ProductDailyRecordInfoExcel
ProductDailyRecordRewardExcel
ProductExcel
ProductMonthlyExcel
ProductSelectExcel
ProductSelectionGroupExcel
GachaCombinedCostExcel
GachaCraftNodeExcel
GachaCraftNodeGroupExcel
GachaCraftOpenTagExcel
GachaGroupExcel
GachaSelectPickupGroupExcel
GachaSelectPickupGroupExcelTable
PickupDuplicateBonusExcel
PickupFirstGetBonus2Excel
PickupFirstGetBonusExcel
DuplicateBonusExcel
```

## 11. 战斗通行证 / BattlePass

```
BattlePassExpLimitExcel
BattlePassFlavorTextExcel
BattlePassInfoExcel
BattlePassLevelExcel
BattlePassMissionExcel
BattlePassRewardExcel
```

## 12. 竞技场扩展 / Arena

```
ArenaLevelSectionExcel
ArenaMapExcel
ArenaNPCExcel               # 注意: gdconf 有 loadArenaNPC,但走 data/ 而非 Excel
ArenaRewardExcel
ConstArenaExcelTable
AssistRewardExcel
AssistSlotExcel
```

## 13. 社交 / SNS / 学生会 / 出勤

```
SNSInfoExcel
SNSPostExcel
SNSProfileExcel
FieldSNSInfoExcel
FieldSNSPostExcel
ClanChattingEmojiExcel
ClanRewardExcel
AttendanceExcel
AttendanceRewardExcel
CafeInteractionExcel
FavorLevelRewardExcel
```

## 14. 活动/迎新活动 / WelcomeCampaign

```
WelcomeCampaignAttendanceRewardExcel
WelcomeCampaignEnterRewardExcel
WelcomeCampaignMissionExcel
WelcomeCampaignSeasonExcel
```

## 15. 限时关卡 / LimitedStage & FarmingDungeon

```
LimitedStageExcelTable
LimitedStageRewardExcelTable
LimitedStageSeasonExcelTable
FarmingDungeonLocationManageExcel
```

## 16. 大场地 / Field（新地图玩法体系）

```
FieldContentStageExcelTable
FieldContentStageRewardExcelTable
FieldCurtainCallFreeModeExcelTable
FieldDateExcelTable
FieldEvidenceExcelTable
FieldInteractionExcelTable
FieldKeywordExcelTable
FieldMasteryExcelTable
FieldMasteryLevelExcelTable
FieldMasteryManageExcelTable
FieldQuestExcelTable
FieldQuestGroupExcel
FieldRewardExcelTable
FieldSceneExcelTable
FieldSeasonExcelTable
FieldStoryStageExcelTable
FieldTutorialExcelTable
FieldWarpExcel
FieldWorldMapZoneExcelTable
ConstFieldExcelTable
```

## 17. 制造 / 配方 / 家具模板

```
RecipeExcel
RecipeSelectionAutoUseExcel
RecipeSelectionGroupExcel
ShiftingCraftRecipeExcel
FurnitureGroupExcel
FurnitureTemplateElementExcel
FurnitureTemplateExcel
DefaultFurnitureExcelTable
ParcelAutoSynthExcel
EquipmentChangePieceExcel
```

## 18. 货币 / 物品扩展 / 默认数据

```
CurrencyExcel
DefaultMailExcelTable
DefaultParcelExcelTable
PresetParcelsExcel
PresetCharacterGroupExcel
PresetCharacterGroupSettingExcel
PossessionCheckExcel
ContentTargetGroupExcel
ContentEnterCostReduceExcel
ContentsFeverExcel
ContentsShortcutExcel
ShortcutTypeExcel
OpenConditionExcel
LevelExpMasterCoinExcel
TagExcelTable
TrophyCollectionExcel
StickerGroupExcel
OperatorExcel
ServiceActionExcel
InformationExcel
CouponStuffExcelTable
ProtocolSettingExcelTable
KeyMappingExcel
KeyMappingPopupExcel
ConstKeyMappingExcelTable
ConstCommonExcelTable
ConstContentsExcelTable
ConstNewbieContentExcelTable
ConstAudioExcelTable
ToastExcel
AlertPopupExcel
MessagePopupExcel
ContentSpoilerPopupExcel
TutorialExcel
TutorialFailureImageExcel
```

## 19. 引导任务扩展 / GuideMission

```
GuideMissionOpenStageConditionExcel
GuideMissionSeasonExcel
```

## 20. Raid 扩展 / TimeAttackDungeon

```
TimeAttackDungeonExcel
PermanentRaidManageExcel
RaidContentPlayGuideExcel
RaidSkillDescriptionListExcel
EliminateRaidStageLimitedRewardExcel
MultiFloorRaidStatChangeExcel
```

## 21. 周常副本扩展 / WeekDungeon

```
WeekDungeonGroupBuffExcel
WeekDungeonOpenScheduleExcel
```

## 22. 账号等级奖励 / AccountLevel

```
AccountLevelRewardExcel
```

## 23. 战役扩展 / Campaign

```
CampaignChapterExcel
CampaignChapterRewardExcel
CampaignStageRewardExcel
```

---

## 备注:命名差异(已实现 vs 未实现的同名变体)

以下表在后端已有"Table版"(来自 Excel.zip),未实现的是"Excel版"(来自 ExcelDB.db),
两者数据重叠,跳过 Excel 版属正常:

| 未实现(跳过)                | 已实现(Loaded)                  |
|------------------------------|----------------------------------|
| CharacterExcel               | CharacterExcelTable              |
| CharacterLevelExcel          | CharacterLevelExcelTable         |
| CharacterStatExcel           | CharacterStatExcelTable          |
| CharacterSkillListExcel      | CharacterSkillListExcelTable     |
| DefaultFurnitureExcelTable   | DefaultFurnitureExcel            |
| GachaSelectPickupGroupExcel  | (仅 Table 版被跳过,均未实现)     |
