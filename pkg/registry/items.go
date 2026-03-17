// Package registry provides Minecraft 1.20.1 item registry.
package registry

import (
	"sync"
)

// ItemID represents an item ID.
type ItemID int32

// ItemProperties holds item properties.
type ItemProperties struct {
	ID            ItemID
	Name          string
	MaxStackSize  int32
	MaxDurability int32
	Edible        bool
	Saturation    float32
	Tool          string
	ToolLevel     int
	AttackDamage  float32
	AttackSpeed   float32
	ArmorMaterial string
	Protection    float32
	CanPlaceOn    []string
	CanDestroy    []string
	BlockID       BlockID
}

// Item represents a stack of items.
type Item struct {
	ID        ItemID
	Count     int32
	Damage    int32 // For items with durability
	NBTData   map[string]interface{}
}

var (
	itemRegistry      map[ItemID]*ItemProperties
	itemRegistryByName map[string]ItemID
	itemMutex         sync.RWMutex
	itemCount         int
)

func init() {
	initializeItems()
}

// Item constants for Minecraft 1.20.1
const (
	ItemAir ItemID = 0
	ItemStone ItemID = 1
	ItemGranite ItemID = 2
	ItemPolishedGranite ItemID = 3
	ItemDiorite ItemID = 4
	ItemPolishedDiorite ItemID = 5
	ItemAndesite ItemID = 6
	ItemPolishedAndesite ItemID = 7
	ItemDeepslate ItemID = 8
	ItemCobbledDeepslate ItemID = 9
	ItemPolishedDeepslate ItemID = 10
	ItemCalcite ItemID = 11
	ItemTuff ItemID = 12
	ItemDirt ItemID = 13
	ItemCoarseDirt ItemID = 14
	ItemPodzol ItemID = 15
	ItemCrimsonNylium ItemID = 16
	ItemWarpedNylium ItemID = 17
	ItemGrassBlock ItemID = 18
	ItemMycelium ItemID = 19
	ItemRootedDirt ItemID = 20
	ItemMossBlock ItemID = 21
	ItemMossCarpet ItemID = 22
	ItemWater ItemID = 23
	ItemLava ItemID = 24
	ItemSand ItemID = 25
	ItemRedSand ItemID = 26
	ItemGravel ItemID = 27
	ItemSandstone ItemID = 28
	ItemChiseledSandstone ItemID = 29
	ItemCutSandstone ItemID = 30
	ItemSmoothSandstone ItemID = 31
	ItemRedSandstone ItemID = 32
	ItemChiseledRedSandstone ItemID = 33
	ItemCutRedSandstone ItemID = 34
	ItemSmoothRedSandstone ItemID = 35
	ItemClay ItemID = 36
	ItemMud ItemID = 37
	ItemPackedMud ItemID = 38
	ItemMudBricks ItemID = 39
	ItemOakLog ItemID = 40
	ItemSpruceLog ItemID = 41
	ItemBirchLog ItemID = 42
	ItemJungleLog ItemID = 43
	ItemAcaciaLog ItemID = 44
	ItemDarkOakLog ItemID = 45
	ItemCrimsonStem ItemID = 46
	ItemWarpedStem ItemID = 47
	ItemStrippedOakLog ItemID = 48
	ItemStrippedSpruceLog ItemID = 49
	ItemStrippedBirchLog ItemID = 50
	ItemStrippedJungleLog ItemID = 51
	ItemStrippedAcaciaLog ItemID = 52
	ItemStrippedDarkOakLog ItemID = 53
	ItemStrippedCrimsonStem ItemID = 54
	ItemStrippedWarpedStem ItemID = 55
	ItemOakWood ItemID = 56
	ItemSpruceWood ItemID = 57
	ItemBirchWood ItemID = 58
	ItemJungleWood ItemID = 59
	ItemAcaciaWood ItemID = 60
	ItemDarkOakWood ItemID = 61
	ItemCrimsonHyphae ItemID = 62
	ItemWarpedHyphae ItemID = 63
	ItemStrippedOakWood ItemID = 64
	ItemStrippedSpruceWood ItemID = 65
	ItemStrippedBirchWood ItemID = 66
	ItemStrippedJungleWood ItemID = 67
	ItemStrippedAcaciaWood ItemID = 68
	ItemStrippedDarkOakWood ItemID = 69
	ItemStrippedCrimsonHyphae ItemID = 70
	ItemStrippedWarpedHyphae ItemID = 71
	ItemOakPlanks ItemID = 72
	ItemSprucePlanks ItemID = 73
	ItemBirchPlanks ItemID = 74
	ItemJunglePlanks ItemID = 75
	ItemAcaciaPlanks ItemID = 76
	ItemDarkOakPlanks ItemID = 77
	ItemCrimsonPlanks ItemID = 78
	ItemWarpedPlanks ItemID = 79
	ItemBambooPlanks ItemID = 80
	ItemMangrovePlanks ItemID = 81
	ItemCherryPlanks ItemID = 82
	ItemBambooMosaic ItemID = 83
	ItemOakSlab ItemID = 84
	ItemSpruceSlab ItemID = 85
	ItemBirchSlab ItemID = 86
	ItemJungleSlab ItemID = 87
	ItemAcaciaSlab ItemID = 88
	ItemDarkOakSlab ItemID = 89
	ItemCrimsonSlab ItemID = 90
	ItemWarpedSlab ItemID = 91
	ItemBambooSlab ItemID = 92
	ItemMangroveSlab ItemID = 93
	ItemCherrySlab ItemID = 94
	ItemBambooMosaicSlab ItemID = 95
	ItemOakStairs ItemID = 96
	ItemSpruceStairs ItemID = 97
	ItemBirchStairs ItemID = 98
	ItemJungleStairs ItemID = 99
	ItemAcaciaStairs ItemID = 100
	ItemDarkOakStairs ItemID = 101
	ItemCrimsonStairs ItemID = 102
	ItemWarpedStairs ItemID = 103
	ItemBambooStairs ItemID = 104
	ItemMangroveStairs ItemID = 105
	ItemCherryStairs ItemID = 106
	ItemGlass ItemID = 107
	ItemTintedGlass ItemID = 108
	ItemWhiteStainedGlass ItemID = 109
	ItemOrangeStainedGlass ItemID = 110
	ItemMagentaStainedGlass ItemID = 111
	ItemLightBlueStainedGlass ItemID = 112
	ItemYellowStainedGlass ItemID = 113
	ItemLimeStainedGlass ItemID = 114
	ItemPinkStainedGlass ItemID = 115
	ItemGrayStainedGlass ItemID = 116
	ItemLightGrayStainedGlass ItemID = 117
	ItemCyanStainedGlass ItemID = 118
	ItemPurpleStainedGlass ItemID = 119
	ItemBlueStainedGlass ItemID = 120
	ItemBrownStainedGlass ItemID = 121
	ItemGreenStainedGlass ItemID = 122
	ItemRedStainedGlass ItemID = 123
	ItemBlackStainedGlass ItemID = 124
	ItemGlassPane ItemID = 125
	ItemWhiteStainedGlassPane ItemID = 126
	ItemOrangeStainedGlassPane ItemID = 127
	ItemMagentaStainedGlassPane ItemID = 128
	ItemLightBlueStainedGlassPane ItemID = 129
	ItemYellowStainedGlassPane ItemID = 130
	ItemLimeStainedGlassPane ItemID = 131
	ItemPinkStainedGlassPane ItemID = 132
	ItemGrayStainedGlassPane ItemID = 133
	ItemLightGrayStainedGlassPane ItemID = 134
	ItemCyanStainedGlassPane ItemID = 135
	ItemPurpleStainedGlassPane ItemID = 136
	ItemBlueStainedGlassPane ItemID = 137
	ItemBrownStainedGlassPane ItemID = 138
	ItemGreenStainedGlassPane ItemID = 139
	ItemRedStainedGlassPane ItemID = 140
	ItemBlackStainedGlassPane ItemID = 141
	ItemPrismarine ItemID = 142
	ItemPrismarineBricks ItemID = 143
	ItemDarkPrismarine ItemID = 144
	ItemSeaLantern ItemID = 145
	ItemCopperBlock ItemID = 146
	ItemExposedCopper ItemID = 147
	ItemWeatheredCopper ItemID = 148
	ItemOxidizedCopper ItemID = 149
	ItemCutCopper ItemID = 150
	ItemExposedCutCopper ItemID = 151
	ItemWeatheredCutCopper ItemID = 152
	ItemOxidizedCutCopper ItemID = 153
	ItemChiseledCopper ItemID = 154
	ItemExposedChiseledCopper ItemID = 155
	ItemWeatheredChiseledCopper ItemID = 156
	ItemOxidizedChiseledCopper ItemID = 157
	ItemCopperGrate ItemID = 158
	ItemExposedCopperGrate ItemID = 159
	ItemWeatheredCopperGrate ItemID = 160
	ItemOxidizedCopperGrate ItemID = 161
	ItemCopperBulb ItemID = 162
	ItemExposedCopperBulb ItemID = 163
	ItemWeatheredCopperBulb ItemID = 164
	ItemOxidizedCopperBulb ItemID = 165
	ItemCopperDoor ItemID = 166
	ItemExposedCopperDoor ItemID = 167
	ItemWeatheredCopperDoor ItemID = 168
	ItemOxidizedCopperDoor ItemID = 169
	ItemCopperTrapdoor ItemID = 170
	ItemExposedCopperTrapdoor ItemID = 171
	ItemWeatheredCopperTrapdoor ItemID = 172
	ItemOxidizedCopperTrapdoor ItemID = 173
	ItemCopperGrateSlab ItemID = 174
	ItemExposedCopperGrateSlab ItemID = 175
	ItemWeatheredCopperGrateSlab ItemID = 176
	ItemOxidizedCopperGrateSlab ItemID = 177
	ItemWaxedCopper ItemID = 178
	ItemWaxedExposedCopper ItemID = 179
	ItemWaxedWeatheredCopper ItemID = 180
	ItemWaxedOxidizedCopper ItemID = 181
	ItemWaxedCutCopper ItemID = 182
	ItemWaxedExposedCutCopper ItemID = 183
	ItemWaxedWeatheredCutCopper ItemID = 184
	ItemWaxedOxidizedCutCopper ItemID = 185
	ItemWaxedChiseledCopper ItemID = 186
	ItemWaxedExposedChiseledCopper ItemID = 187
	ItemWaxedWeatheredChiseledCopper ItemID = 188
	ItemWaxedOxidizedChiseledCopper ItemID = 189
	ItemWaxedCopperGrate ItemID = 190
	ItemWaxedExposedCopperGrate ItemID = 191
	ItemWaxedWeatheredCopperGrate ItemID = 192
	ItemWaxedOxidizedCopperGrate ItemID = 193
	ItemWaxedCopperBulb ItemID = 194
	ItemWaxedExposedCopperBulb ItemID = 195
	ItemWaxedWeatheredCopperBulb ItemID = 196
	ItemWaxedOxidizedCopperBulb ItemID = 197
	ItemWaxedCopperDoor ItemID = 198
	ItemWaxedExposedCopperDoor ItemID = 199
	ItemWaxedWeatheredCopperDoor ItemID = 200
	ItemWaxedOxidizedCopperDoor ItemID = 201
	ItemWaxedCopperTrapdoor ItemID = 202
	ItemWaxedExposedCopperTrapdoor ItemID = 203
	ItemWaxedWeatheredCopperTrapdoor ItemID = 204
	ItemWaxedOxidizedCopperTrapdoor ItemID = 205
	ItemWaxedCopperGrateSlab ItemID = 206
	ItemWaxedExposedCopperGrateSlab ItemID = 207
	ItemWaxedWeatheredCopperGrateSlab ItemID = 208
	ItemWaxedOxidizedCopperGrateSlab ItemID = 209
	ItemIronBlock ItemID = 210
	ItemCauldron ItemID = 211
	ItemWaterCauldron ItemID = 212
	ItemLavaCauldron ItemID = 213
	ItemPowderSnowCauldron ItemID = 214
	ItemCopperSulfateCauldron ItemID = 215
	ItemGoldBlock ItemID = 216
	ItemRawIronBlock ItemID = 217
	ItemRawCopperBlock ItemID = 218
	ItemRawGoldBlock ItemID = 219
	ItemAmethystBlock ItemID = 220
	ItemBuddingAmethyst ItemID = 221
	ItemCopperOre ItemID = 222
	ItemDeepslateCopperOre ItemID = 223
	ItemIronOre ItemID = 224
	ItemDeepslateIronOre ItemID = 225
	ItemGoldOre ItemID = 226
	ItemDeepslateGoldOre ItemID = 227
	ItemDiamondOre ItemID = 228
	ItemDeepslateDiamondOre ItemID = 229
	ItemLapisOre ItemID = 230
	ItemDeepslateLapisOre ItemID = 231
	ItemRedstoneOre ItemID = 232
	ItemDeepslateRedstoneOre ItemID = 233
	ItemEmeraldOre ItemID = 234
	ItemDeepslateEmeraldOre ItemID = 235
	ItemCoalOre ItemID = 236
	ItemDeepslateCoalOre ItemID = 237
	ItemNetherQuartzOre ItemID = 238
	ItemNetherGoldOre ItemID = 239
	ItemAncientDebris ItemID = 240
	ItemCoalBlock ItemID = 241
	ItemDiamondBlock ItemID = 242
	ItemLapisBlock ItemID = 243
	ItemEmeraldBlock ItemID = 244
	ItemRedstoneBlock ItemID = 245
	ItemNetheriteBlock ItemID = 246
	ItemOakFence ItemID = 247
	ItemSpruceFence ItemID = 248
	ItemBirchFence ItemID = 249
	ItemJungleFence ItemID = 250
	ItemAcaciaFence ItemID = 251
	ItemDarkOakFence ItemID = 252
	ItemCrimsonFence ItemID = 253
	ItemWarpedFence ItemID = 254
	ItemMangroveFence ItemID = 255
	ItemCherryFence ItemID = 256
	ItemBambooFence ItemID = 257
	ItemOakFenceGate ItemID = 258
	ItemSpruceFenceGate ItemID = 259
	ItemBirchFenceGate ItemID = 260
	ItemJungleFenceGate ItemID = 261
	ItemAcaciaFenceGate ItemID = 262
	ItemDarkOakFenceGate ItemID = 263
	ItemCrimsonFenceGate ItemID = 264
	ItemWarpedFenceGate ItemID = 265
	ItemMangroveFenceGate ItemID = 266
	ItemCherryFenceGate ItemID = 267
	ItemBambooFenceGate ItemID = 268
	ItemOakDoor ItemID = 269
	ItemSpruceDoor ItemID = 270
	ItemBirchDoor ItemID = 271
	ItemJungleDoor ItemID = 272
	ItemAcaciaDoor ItemID = 273
	ItemDarkOakDoor ItemID = 274
	ItemCrimsonDoor ItemID = 275
	ItemWarpedDoor ItemID = 276
	ItemMangroveDoor ItemID = 277
	ItemCherryDoor ItemID = 278
	ItemBambooDoor ItemID = 279
	ItemIronDoor ItemID = 280
	ItemOakTrapdoor ItemID = 289
	ItemSpruceTrapdoor ItemID = 290
	ItemBirchTrapdoor ItemID = 291
	ItemJungleTrapdoor ItemID = 292
	ItemAcaciaTrapdoor ItemID = 293
	ItemDarkOakTrapdoor ItemID = 294
	ItemCrimsonTrapdoor ItemID = 295
	ItemWarpedTrapdoor ItemID = 296
	ItemMangroveTrapdoor ItemID = 297
	ItemCherryTrapdoor ItemID = 298
	ItemBambooTrapdoor ItemID = 299
	ItemIronTrapdoor ItemID = 300
	ItemOakPressurePlate ItemID = 309
	ItemSprucePressurePlate ItemID = 310
	ItemBirchPressurePlate ItemID = 311
	ItemJunglePressurePlate ItemID = 312
	ItemAcaciaPressurePlate ItemID = 313
	ItemDarkOakPressurePlate ItemID = 314
	ItemCrimsonPressurePlate ItemID = 315
	ItemWarpedPressurePlate ItemID = 316
	ItemMangrovePressurePlate ItemID = 317
	ItemCherryPressurePlate ItemID = 318
	ItemBambooPressurePlate ItemID = 319
	ItemStonePressurePlate ItemID = 320
	ItemPolishedBlackstonePressurePlate ItemID = 321
	ItemHeavyWeightedPressurePlate ItemID = 322
	ItemLightWeightedPressurePlate ItemID = 323
	ItemOakButton ItemID = 332
	ItemSpruceButton ItemID = 333
	ItemBirchButton ItemID = 334
	ItemJungleButton ItemID = 335
	ItemAcaciaButton ItemID = 336
	ItemDarkOakButton ItemID = 337
	ItemCrimsonButton ItemID = 338
	ItemWarpedButton ItemID = 339
	ItemMangroveButton ItemID = 340
	ItemCherryButton ItemID = 341
	ItemBambooButton ItemID = 342
	ItemStoneButton ItemID = 343
	ItemPolishedBlackstoneButton ItemID = 344
	ItemBricks ItemID = 353
	ItemBrickSlab ItemID = 354
	ItemBrickStairs ItemID = 355
	ItemBrickWall ItemID = 356
	ItemMudBrickSlab ItemID = 357
	ItemMudBrickStairs ItemID = 358
	ItemMudBrickWall ItemID = 359
	ItemCobbledDeepslateSlab ItemID = 360
	ItemCobbledDeepslateStairs ItemID = 361
	ItemCobbledDeepslateWall ItemID = 362
	ItemPolishedDeepslateSlab ItemID = 363
	ItemPolishedDeepslateStairs ItemID = 364
	ItemPolishedDeepslateWall ItemID = 365
	ItemDeepslateBricks ItemID = 366
	ItemDeepslateBrickSlab ItemID = 367
	ItemDeepslateBrickStairs ItemID = 368
	ItemDeepslateBrickWall ItemID = 369
	ItemDeepslateTiles ItemID = 370
	ItemDeepslateTileSlab ItemID = 371
	ItemDeepslateTileStairs ItemID = 372
	ItemDeepslateTileWall ItemID = 373
	ItemReinforcedDeepslate ItemID = 374
	ItemCobblestone ItemID = 375
	ItemCobblestoneSlab ItemID = 376
	ItemCobblestoneStairs ItemID = 377
	ItemCobblestoneWall ItemID = 378
	ItemMossyCobblestone ItemID = 379
	ItemMossyCobblestoneSlab ItemID = 380
	ItemMossyCobblestoneStairs ItemID = 381
	ItemMossyCobblestoneWall ItemID = 382
	ItemStoneBricks ItemID = 383
	ItemStoneBrickSlab ItemID = 384
	ItemStoneBrickStairs ItemID = 385
	ItemStoneBrickWall ItemID = 386
	ItemMossyStoneBricks ItemID = 387
	ItemCrackedStoneBricks ItemID = 388
	ItemChiseledStoneBricks ItemID = 389
	ItemSmoothStone ItemID = 390
	ItemSmoothStoneSlab ItemID = 391
	ItemSandstoneSlab ItemID = 392
	ItemSandstoneStairs ItemID = 393
	ItemSandstoneWall ItemID = 394
	ItemRedSandstoneSlab ItemID = 395
	ItemRedSandstoneStairs ItemID = 396
	ItemRedSandstoneWall ItemID = 397
	ItemPrismarineSlab ItemID = 398
	ItemPrismarineStairs ItemID = 399
	ItemPrismarineWall ItemID = 400
	ItemPrismarineBrickSlab ItemID = 401
	ItemPrismarineBrickStairs ItemID = 402
	ItemDarkPrismarineSlab ItemID = 403
	ItemDarkPrismarineStairs ItemID = 404
	ItemNetherBricks ItemID = 405
	ItemNetherBrickSlab ItemID = 406
	ItemNetherBrickStairs ItemID = 407
	ItemNetherBrickFence ItemID = 408
	ItemNetherBrickWall ItemID = 409
	ItemRedNetherBricks ItemID = 410
	ItemBlackstone ItemID = 411
	ItemBlackstoneSlab ItemID = 412
	ItemBlackstoneStairs ItemID = 413
	ItemBlackstoneWall ItemID = 414
	ItemPolishedBlackstone ItemID = 415
	ItemPolishedBlackstoneSlab ItemID = 416
	ItemPolishedBlackstoneStairs ItemID = 417
	ItemPolishedBlackstoneWall ItemID = 418
	ItemPolishedBlackstoneBricks ItemID = 419
	ItemPolishedBlackstoneBrickSlab ItemID = 420
	ItemPolishedBlackstoneBrickStairs ItemID = 421
	ItemPolishedBlackstoneBrickWall ItemID = 422
	ItemChiseledPolishedBlackstone ItemID = 423
	ItemCrackedPolishedBlackstoneBricks ItemID = 424
	ItemGildedBlackstone ItemID = 425
	ItemEndStoneBricks ItemID = 426
	ItemEndStoneBrickSlab ItemID = 427
	ItemEndStoneBrickStairs ItemID = 428
	ItemEndStoneBrickWall ItemID = 429
	ItemPurpurBlock ItemID = 430
	ItemPurpurPillar ItemID = 431
	ItemPurpurSlab ItemID = 432
	ItemPurpurStairs ItemID = 433
	ItemQuartzBlock ItemID = 434
	ItemChiseledQuartzBlock ItemID = 435
	ItemQuartzPillar ItemID = 436
	ItemQuartzBricks ItemID = 437
	ItemQuartzSlab ItemID = 438
	ItemQuartzStairs ItemID = 439
	ItemSmoothQuartz ItemID = 440
	ItemTerracotta ItemID = 441
	ItemWhiteTerracotta ItemID = 442
	ItemOrangeTerracotta ItemID = 443
	ItemMagentaTerracotta ItemID = 444
	ItemLightBlueTerracotta ItemID = 445
	ItemYellowTerracotta ItemID = 446
	ItemLimeTerracotta ItemID = 447
	ItemPinkTerracotta ItemID = 448
	ItemGrayTerracotta ItemID = 449
	ItemLightGrayTerracotta ItemID = 450
	ItemCyanTerracotta ItemID = 451
	ItemPurpleTerracotta ItemID = 452
	ItemBlueTerracotta ItemID = 453
	ItemBrownTerracotta ItemID = 454
	ItemGreenTerracotta ItemID = 455
	ItemRedTerracotta ItemID = 456
	ItemBlackTerracotta ItemID = 457
	ItemWhiteGlazedTerracotta ItemID = 458
	ItemOrangeGlazedTerracotta ItemID = 459
	ItemMagentaGlazedTerracotta ItemID = 460
	ItemLightBlueGlazedTerracotta ItemID = 461
	ItemYellowGlazedTerracotta ItemID = 462
	ItemLimeGlazedTerracotta ItemID = 463
	ItemPinkGlazedTerracotta ItemID = 464
	ItemGrayGlazedTerracotta ItemID = 465
	ItemLightGrayGlazedTerracotta ItemID = 466
	ItemCyanGlazedTerracotta ItemID = 467
	ItemPurpleGlazedTerracotta ItemID = 468
	ItemBlueGlazedTerracotta ItemID = 469
	ItemBrownGlazedTerracotta ItemID = 470
	ItemGreenGlazedTerracotta ItemID = 471
	ItemRedGlazedTerracotta ItemID = 472
	ItemBlackGlazedTerracotta ItemID = 473
	ItemWhiteWool ItemID = 474
	ItemOrangeWool ItemID = 475
	ItemMagentaWool ItemID = 476
	ItemLightBlueWool ItemID = 477
	ItemYellowWool ItemID = 478
	ItemLimeWool ItemID = 479
	ItemPinkWool ItemID = 480
	ItemGrayWool ItemID = 481
	ItemLightGrayWool ItemID = 482
	ItemCyanWool ItemID = 483
	ItemPurpleWool ItemID = 484
	ItemBlueWool ItemID = 485
	ItemBrownWool ItemID = 486
	ItemGreenWool ItemID = 487
	ItemRedWool ItemID = 488
	ItemBlackWool ItemID = 489
	ItemWhiteCarpet ItemID = 490
	ItemOrangeCarpet ItemID = 491
	ItemMagentaCarpet ItemID = 492
	ItemLightBlueCarpet ItemID = 493
	ItemYellowCarpet ItemID = 494
	ItemLimeCarpet ItemID = 495
	ItemPinkCarpet ItemID = 496
	ItemGrayCarpet ItemID = 497
	ItemLightGrayCarpet ItemID = 498
	ItemCyanCarpet ItemID = 499
	ItemPurpleCarpet ItemID = 500
	ItemBlueCarpet ItemID = 501
	ItemBrownCarpet ItemID = 502
	ItemGreenCarpet ItemID = 503
	ItemRedCarpet ItemID = 504
	ItemBlackCarpet ItemID = 505
	ItemCraftingTable ItemID = 506
	ItemFurnace ItemID = 507
	ItemSmoker ItemID = 508
	ItemBlastFurnace ItemID = 509
	ItemLectern ItemID = 510
	ItemLoom ItemID = 511
	ItemCartographyTable ItemID = 512
	ItemFletchingTable ItemID = 513
	ItemGrindstone ItemID = 514
	ItemSmithingTable ItemID = 515
	ItemStonecutter ItemID = 516
	ItemBrewingStand ItemID = 517
	ItemBarrel ItemID = 518
	ItemChiseledBookshelf ItemID = 519
	ItemDecoratedPot ItemID = 520
	ItemShulkerBox ItemID = 521
	ItemWhiteShulkerBox ItemID = 522
	ItemOrangeShulkerBox ItemID = 523
	ItemMagentaShulkerBox ItemID = 524
	ItemLightBlueShulkerBox ItemID = 525
	ItemYellowShulkerBox ItemID = 526
	ItemLimeShulkerBox ItemID = 527
	ItemPinkShulkerBox ItemID = 528
	ItemGrayShulkerBox ItemID = 529
	ItemLightGrayShulkerBox ItemID = 530
	ItemCyanShulkerBox ItemID = 531
	ItemPurpleShulkerBox ItemID = 532
	ItemBlueShulkerBox ItemID = 533
	ItemBrownShulkerBox ItemID = 534
	ItemGreenShulkerBox ItemID = 535
	ItemRedShulkerBox ItemID = 536
	ItemBlackShulkerBox ItemID = 537
	ItemChest ItemID = 538
	ItemTrappedChest ItemID = 539
	ItemEnderChest ItemID = 540
	ItemBedrock ItemID = 564
	ItemWhiteConcrete ItemID = 566
	ItemOrangeConcrete ItemID = 567
	ItemMagentaConcrete ItemID = 568
	ItemLightBlueConcrete ItemID = 569
	ItemYellowConcrete ItemID = 570
	ItemLimeConcrete ItemID = 571
	ItemPinkConcrete ItemID = 572
	ItemGrayConcrete ItemID = 573
	ItemLightGrayConcrete ItemID = 574
	ItemCyanConcrete ItemID = 575
	ItemPurpleConcrete ItemID = 576
	ItemBlueConcrete ItemID = 577
	ItemBrownConcrete ItemID = 578
	ItemGreenConcrete ItemID = 579
	ItemRedConcrete ItemID = 580
	ItemBlackConcrete ItemID = 581
	ItemWhiteConcretePowder ItemID = 582
	ItemOrangeConcretePowder ItemID = 583
	ItemMagentaConcretePowder ItemID = 584
	ItemLightBlueConcretePowder ItemID = 585
	ItemYellowConcretePowder ItemID = 586
	ItemLimeConcretePowder ItemID = 587
	ItemPinkConcretePowder ItemID = 588
	ItemGrayConcretePowder ItemID = 589
	ItemLightGrayConcretePowder ItemID = 590
	ItemCyanConcretePowder ItemID = 591
	ItemPurpleConcretePowder ItemID = 592
	ItemBlueConcretePowder ItemID = 593
	ItemBrownConcretePowder ItemID = 594
	ItemGreenConcretePowder ItemID = 595
	ItemRedConcretePowder ItemID = 596
	ItemBlackConcretePowder ItemID = 597
	ItemSponge ItemID = 598
	ItemWetSponge ItemID = 599
	ItemGlowstone ItemID = 600
	ItemSeaPickle ItemID = 601
	ItemSnowBlock ItemID = 602
	ItemSnow ItemID = 603
	ItemIce ItemID = 604
	ItemPackedIce ItemID = 605
	ItemBlueIce ItemID = 606
	ItemPowderSnow ItemID = 607
	ItemFire ItemID = 608
	ItemSoulFire ItemID = 609
	ItemCampfire ItemID = 610
	ItemSoulCampfire ItemID = 611
	ItemTorch ItemID = 612
	ItemSoulTorch ItemID = 613
	ItemEndRod ItemID = 614
	ItemLantern ItemID = 615
	ItemSoulLantern ItemID = 616
	ItemGlowLichen ItemID = 617
	ItemLight ItemID = 618
	ItemLightningRod ItemID = 619
	ItemSculk ItemID = 629
	ItemSculkSensor ItemID = 630
	ItemSculkCatalyst ItemID = 631
	ItemSculkVein ItemID = 632
	ItemSculkShrieker ItemID = 633
	ItemNetherPortal ItemID = 635
	ItemEndPortal ItemID = 636
	ItemEndPortalFrame ItemID = 637
	ItemEndGateway ItemID = 638
	ItemCommandBlock ItemID = 639
	ItemChainCommandBlock ItemID = 640
	ItemRepeatingCommandBlock ItemID = 641
	ItemStructureBlock ItemID = 642
	ItemJigsaw ItemID = 643
	ItemSpawner ItemID = 644
	ItemMonsterEgg ItemID = 645
	ItemObsidian ItemID = 646
	ItemCryingObsidian ItemID = 647
	ItemDragonEgg ItemID = 648
	ItemNetherStar ItemID = 649
	ItemEnchantingTable ItemID = 650
	ItemBookshelf ItemID = 651
	ItemCocoa ItemID = 661
	ItemWheat ItemID = 662
	ItemCarrots ItemID = 663
	ItemPotatoes ItemID = 664
	ItemBeetroots ItemID = 665
	ItemMelonStem ItemID = 666
	ItemPumpkinStem ItemID = 667
	ItemSweetBerryBush ItemID = 668
	ItemVine ItemID = 669
	ItemGlowBerryBlock ItemID = 670
	ItemCaveVines ItemID = 671
	ItemCaveVinesBody ItemID = 672
	ItemTwistingVines ItemID = 673
	ItemTwistingVinesPlant ItemID = 674
	ItemWeepingVines ItemID = 675
	ItemWeepingVinesPlant ItemID = 676
	ItemLilyPad ItemID = 677
	ItemSunflower ItemID = 678
	ItemLilac ItemID = 679
	ItemRoseBush ItemID = 680
	ItemPeony ItemID = 681
	ItemTallGrass ItemID = 682
	ItemLargeFern ItemID = 683
	ItemPinkPetals ItemID = 684
	ItemMangrovePropagule ItemID = 685
	ItemSapling ItemID = 687
	ItemOakSapling ItemID = 688
	ItemSpruceSapling ItemID = 689
	ItemBirchSapling ItemID = 690
	ItemJungleSapling ItemID = 691
	ItemAcaciaSapling ItemID = 692
	ItemDarkOakSapling ItemID = 693
	ItemCherrySapling ItemID = 694
	ItemMangrovePropaguleBlock ItemID = 695
	ItemAzalea ItemID = 696
	ItemFloweringAzalea ItemID = 697
	ItemSporeBlossom ItemID = 698
	ItemGrass ItemID = 699
	ItemFern ItemID = 700
	ItemDeadBush ItemID = 701
	ItemSeagrass ItemID = 702
	ItemTallSeagrass ItemID = 703
	ItemKelp ItemID = 704
	ItemKelpPlant ItemID = 705
	ItemBamboo ItemID = 706
	ItemBambooSapling ItemID = 707
	ItemBambooShoot ItemID = 708
	ItemCactus ItemID = 709
	ItemSugarCane ItemID = 710
	ItemCoral ItemID = 711
	ItemCoralBlock ItemID = 712
	ItemCoralFan ItemID = 713
	ItemCoralWallFan ItemID = 714
	ItemBrainCoral ItemID = 715
	ItemBubbleCoral ItemID = 716
	ItemFireCoral ItemID = 717
	ItemHornCoral ItemID = 718
	ItemTubeCoral ItemID = 719
	ItemBrainCoralBlock ItemID = 720
	ItemBubbleCoralBlock ItemID = 721
	ItemFireCoralBlock ItemID = 722
	ItemHornCoralBlock ItemID = 723
	ItemTubeCoralBlock ItemID = 724
	ItemBrainCoralFan ItemID = 725
	ItemBubbleCoralFan ItemID = 726
	ItemFireCoralFan ItemID = 727
	ItemHornCoralFan ItemID = 728
	ItemTubeCoralFan ItemID = 729
	ItemBrainCoralWallFan ItemID = 730
	ItemBubbleCoralWallFan ItemID = 731
	ItemFireCoralWallFan ItemID = 732
	ItemHornCoralWallFan ItemID = 733
	ItemTubeCoralWallFan ItemID = 734
	ItemCoralPlant ItemID = 735
	ItemChorusPlant ItemID = 736
	ItemChorusFlower ItemID = 737
	ItemCrimsonFungus ItemID = 738
	ItemWarpedFungus ItemID = 739
	ItemCrimsonRoots ItemID = 740
	ItemWarpedRoots ItemID = 741
	ItemNetherSprouts ItemID = 742
	ItemWeepingVinesBlock ItemID = 743
	ItemTwistingVinesBlock ItemID = 744
	ItemBigDripleaf ItemID = 745
	ItemBigDripleafStem ItemID = 746
	ItemSmallDripleaf ItemID = 747
	ItemSporeBlossomBlock ItemID = 748
	ItemMossCarpetBlock ItemID = 749
	ItemPinkPetalsBlock ItemID = 750
	ItemMangroveLeaves ItemID = 751
	ItemMangroveRoots ItemID = 752
	ItemMuddyMangroveRoots ItemID = 754
	ItemCherryLeaves ItemID = 755
	ItemCherryLog ItemID = 756
	ItemCherrySaplingBlock ItemID = 757
	ItemHangingRoots ItemID = 759
	ItemNetherGoldBlock ItemID = 760
	ItemPolishedBasalt ItemID = 761
	ItemBasalt ItemID = 762
	ItemSmoothBasalt ItemID = 763
	ItemSoulSoil ItemID = 766
	ItemSoulSand ItemID = 765
	ItemGrassPath ItemID = 768
	ItemFarmland ItemID = 769
	ItemNetherrack ItemID = 770
	ItemCrimsonNyliumBlock ItemID = 771
	ItemWarpedNyliumBlock ItemID = 772
	ItemMagmaBlock ItemID = 773
	ItemNetherBricksBlock ItemID = 774
	ItemRedNetherBricksBlock ItemID = 775
	ItemNetherQuartzOreBlock ItemID = 776
	ItemNetherGoldOreBlock ItemID = 777
	ItemAncientDebrisBlock ItemID = 778
	ItemBasaltBlock ItemID = 779
	ItemPolishedBasaltBlock ItemID = 780
	ItemSmoothBasaltBlock ItemID = 781
	ItemBlueIceBlock ItemID = 782
	ItemPackedIceBlock ItemID = 783
	ItemIceBlock ItemID = 784
	ItemSnowBlockBlock ItemID = 785
	ItemLadder ItemID = 787
	ItemVines ItemID = 788
	ItemTripwire ItemID = 789
	ItemTripwireHook ItemID = 790
	ItemCobweb ItemID = 791
	ItemHayBlock ItemID = 792
	ItemTarget ItemID = 793
	ItemNoteBlock ItemID = 794
	ItemRedstoneWire ItemID = 795
	ItemRedstoneTorch ItemID = 796
	ItemRedstoneWallTorch ItemID = 797
	ItemRepeater ItemID = 798
	ItemRedstoneRepeater ItemID = 799
	ItemRedstoneComparator ItemID = 800
	ItemComparator ItemID = 801
	ItemPiston ItemID = 802
	ItemStickyPiston ItemID = 803
	ItemPistonHead ItemID = 804
	ItemMovingPiston ItemID = 805
	ItemDispenser ItemID = 806
	ItemDropper ItemID = 807
	ItemObserver ItemID = 808
	ItemHopper ItemID = 809
	ItemLever ItemID = 810
	ItemPoweredRail ItemID = 811
	ItemDetectorRail ItemID = 812
	ItemRail ItemID = 813
	ItemActivatorRail ItemID = 814
	ItemRailBlock ItemID = 815
	ItemPoweredRailBlock ItemID = 816
	ItemDetectorRailBlock ItemID = 817
	ItemRailBlockBlock ItemID = 818
	ItemActivatorRailBlock ItemID = 819
	ItemIronPickaxe ItemID = 900
	ItemIronShovel ItemID = 901
	ItemIronAxe ItemID = 902
	ItemIronSword ItemID = 903
	ItemIronHoe ItemID = 904
	ItemDiamondPickaxe ItemID = 905
	ItemDiamondShovel ItemID = 906
	ItemDiamondAxe ItemID = 907
	ItemDiamondSword ItemID = 908
	ItemDiamondHoe ItemID = 909
	ItemStonePickaxe ItemID = 910
	ItemStoneShovel ItemID = 911
	ItemStoneAxe ItemID = 912
	ItemStoneSword ItemID = 913
	ItemStoneHoe ItemID = 914
	ItemWoodenPickaxe ItemID = 915
	ItemWoodenShovel ItemID = 916
	ItemWoodenAxe ItemID = 917
	ItemWoodenSword ItemID = 918
	ItemWoodenHoe ItemID = 919
	ItemGoldenPickaxe ItemID = 920
	ItemGoldenShovel ItemID = 921
	ItemGoldenAxe ItemID = 922
	ItemGoldenSword ItemID = 923
	ItemGoldenHoe ItemID = 924
	ItemNetheritePickaxe ItemID = 925
	ItemNetheriteShovel ItemID = 926
	ItemNetheriteAxe ItemID = 927
	ItemNetheriteSword ItemID = 928
	ItemNetheriteHoe ItemID = 929
	ItemShears ItemID = 930
	ItemFlintAndSteel ItemID = 931
	ItemShield ItemID = 932
	ItemTrident ItemID = 933
	ItemElytra ItemID = 934
	ItemFishingRod ItemID = 935
	ItemCarrotOnAStick ItemID = 936
	ItemWarpedFungusOnAStick ItemID = 937
	ItemTotemOfUndying ItemID = 938
	ItemEnderPearl ItemID = 939
	ItemArmorStand ItemID = 940
	ItemItemFrame ItemID = 941
	ItemGlowItemFrame ItemID = 942
	ItemPainting ItemID = 943
	ItemFlowerPot ItemID = 944
	ItemWhiteBed ItemID = 945
	ItemOrangeBed ItemID = 946
	ItemMagentaBed ItemID = 947
	ItemLightBlueBed ItemID = 948
	ItemYellowBed ItemID = 949
	ItemLimeBed ItemID = 950
	ItemPinkBed ItemID = 951
	ItemGrayBed ItemID = 952
	ItemLightGrayBed ItemID = 953
	ItemCyanBed ItemID = 954
	ItemPurpleBed ItemID = 955
	ItemBlueBed ItemID = 956
	ItemBrownBed ItemID = 957
	ItemGreenBed ItemID = 958
	ItemRedBed ItemID = 959
	ItemBlackBed ItemID = 960
	ItemLeatherHelmet ItemID = 961
	ItemLeatherChestplate ItemID = 962
	ItemLeatherLeggings ItemID = 963
	ItemLeatherBoots ItemID = 964
	ItemChainmailHelmet ItemID = 965
	ItemChainmailChestplate ItemID = 966
	ItemChainmailLeggings ItemID = 967
	ItemChainmailBoots ItemID = 968
	ItemIronHelmet ItemID = 969
	ItemIronChestplate ItemID = 970
	ItemIronLeggings ItemID = 971
	ItemIronBoots ItemID = 972
	ItemDiamondHelmet ItemID = 973
	ItemDiamondChestplate ItemID = 974
	ItemDiamondLeggings ItemID = 975
	ItemDiamondBoots ItemID = 976
	ItemGoldenHelmet ItemID = 977
	ItemGoldenChestplate ItemID = 978
	ItemGoldenLeggings ItemID = 979
	ItemGoldenBoots ItemID = 980
	ItemNetheriteHelmet ItemID = 981
	ItemNetheriteChestplate ItemID = 982
	ItemNetheriteLeggings ItemID = 983
	ItemNetheriteBoots ItemID = 984
	ItemTurtleHelmet ItemID = 985
	ItemApple ItemID = 986
	ItemGoldenApple ItemID = 987
	ItemEnchantedGoldenApple ItemID = 988
	ItemBread ItemID = 989
	ItemPorkchop ItemID = 990
	ItemCookedPorkchop ItemID = 991
	ItemBeef ItemID = 992
	ItemCookedBeef ItemID = 993
	ItemChicken ItemID = 994
	ItemCookedChicken ItemID = 995
	ItemCod ItemID = 996
	ItemCookedCod ItemID = 997
	ItemSalmon ItemID = 998
	ItemCookedSalmon ItemID = 999
	ItemTropicalFish ItemID = 1000
	ItemPufferfish ItemID = 1001
	ItemMutton ItemID = 1002
	ItemCookedMutton ItemID = 1003
	ItemRabbit ItemID = 1004
	ItemCookedRabbit ItemID = 1005
	ItemRabbitStew ItemID = 1006
	ItemMushroomStew ItemID = 1007
	ItemBeetrootSoup ItemID = 1008
	ItemCookie ItemID = 1009
	ItemCake ItemID = 1010
	ItemMelonSlice ItemID = 1011
	ItemSweetBerries ItemID = 1012
	ItemGlowBerries ItemID = 1013
	ItemPumpkinPie ItemID = 1014
	ItemCarrot ItemID = 1015
	ItemPotato ItemID = 1016
	ItemBakedPotato ItemID = 1017
	ItemPoisonousPotato ItemID = 1018
	ItemGoldenCarrot ItemID = 1019
	ItemBeetroot ItemID = 1020
	ItemDriedKelp ItemID = 1021
	ItemChorusFruit ItemID = 1022
	ItemHoneyBottle ItemID = 1023
	ItemBow ItemID = 1024
	ItemArrow ItemID = 1025
	ItemSpectralArrow ItemID = 1026
	ItemTippedArrow ItemID = 1027
	ItemPotion ItemID = 1028
	ItemSplashPotion ItemID = 1029
	ItemLingeringPotion ItemID = 1030
	ItemCrossbow ItemID = 1031
	ItemFireworkRocket ItemID = 1032
	ItemFireworkStar ItemID = 1033
	ItemWritableBook ItemID = 1034
	ItemWrittenBook ItemID = 1035
	ItemEmerald ItemID = 1036
	ItemDiamond ItemID = 1037
	ItemIronIngot ItemID = 1038
	ItemGoldIngot ItemID = 1039
	ItemCopperIngot ItemID = 1040
	ItemNetheriteIngot ItemID = 1041
	ItemCoal ItemID = 1042
	ItemCharcoal ItemID = 1043
	ItemLapisLazuli ItemID = 1044
	ItemRedstone ItemID = 1045
	ItemQuartz ItemID = 1046
	ItemAmethystShard ItemID = 1047
	ItemEnderPearlItem ItemID = 1048
	ItemGlowstoneDust ItemID = 1049
	ItemBlazeRod ItemID = 1050
	ItemNetherStarItem ItemID = 1051
	ItemDragonBreathItem ItemID = 1052
	ItemString ItemID = 1053
	ItemFeather ItemID = 1054
	ItemStick ItemID = 1055
	ItemBone ItemID = 1056
	ItemBoneMeal ItemID = 1057
	ItemGunpowder ItemID = 1058
	ItemSpiderEye ItemID = 1059
	ItemRottenFlesh ItemID = 1060
	ItemSlimeball ItemID = 1061
	ItemEgg ItemID = 1062
	ItemCompass ItemID = 1063
	ItemRecoveryCompass ItemID = 1064
	ItemClock ItemID = 1065
	ItemSpyglass ItemID = 1066
	ItemLead ItemID = 1067
	ItemNameTag ItemID = 1068
	ItemSaddle ItemID = 1069
	ItemSign ItemID = 1070
	ItemHangingSign ItemID = 1071
	ItemBucket ItemID = 1072
	ItemWaterBucket ItemID = 1073
	ItemLavaBucket ItemID = 1074
	ItemPowderSnowBucket ItemID = 1075
	ItemMilkBucket ItemID = 1076
	ItemCodBucket ItemID = 1077
	ItemSalmonBucket ItemID = 1078
	ItemTropicalFishBucket ItemID = 1079
	ItemPufferfishBucket ItemID = 1080
	ItemAxolotlBucket ItemID = 1081
	ItemTadpoleBucket ItemID = 1082
	ItemGoatHorn ItemID = 1083
	ItemOminousBanner ItemID = 1084
	ItemDiscFragment5 ItemID = 1085
	ItemDiscFragment ItemID = 1086
	ItemPattern ItemID = 1087
	ItemNetheriteUpgrade ItemID = 1088
	ItemSmithingTemplate ItemID = 1089
	ItemDiamondHorseArmor ItemID = 1090
	ItemGoldHorseArmor ItemID = 1091
	ItemIronHorseArmor ItemID = 1092
	ItemLeatherHorseArmor ItemID = 1093
	ItemEnchantedBook ItemID = 1094
	ItemBottleOfEnchanting ItemID = 1095
	ItemHeartOfTheSea ItemID = 1096
	ItemNautilusShell ItemID = 1097
	ItemShulkerShell ItemID = 1098
	ItemBannerPattern ItemID = 1099
	ItemDiscFragment13 ItemID = 1100
	ItemDiscFragmentItem ItemID = 1101
	ItemMusicDisc13 ItemID = 1102
	ItemMusicDiscCat ItemID = 1103
	ItemMusicDiscBlocks ItemID = 1104
	ItemMusicDiscChirp ItemID = 1105
	ItemMusicDiscFar ItemID = 1106
	ItemMusicDiscMall ItemID = 1107
	ItemMusicDiscMellohi ItemID = 1108
	ItemMusicDiscStal ItemID = 1109
	ItemMusicDiscStrad ItemID = 1110
	ItemMusicDiscWard ItemID = 1111
	ItemMusicDisc11 ItemID = 1112
	ItemMusicDiscWait ItemID = 1113
	ItemMusicDiscOtherside ItemID = 1114
	ItemMusicDisc5 ItemID = 1115
	ItemMusicDiscPigstep ItemID = 1116
	ItemMusicDiscRelic ItemID = 1117
	ItemMusicDiscDiscFragment ItemID = 1118
	ItemMusicDiscCreator ItemID = 1119
	ItemMusicDiscCreatorMusicBox ItemID = 1120
	ItemGlowInkSac ItemID = 1121
	ItemInkSac ItemID = 1122
	ItemCocoaBeans ItemID = 1123
	ItemWhiteDye ItemID = 1124
	ItemOrangeDye ItemID = 1125
	ItemMagentaDye ItemID = 1126
	ItemLightBlueDye ItemID = 1127
	ItemYellowDye ItemID = 1128
	ItemLimeDye ItemID = 1129
	ItemPinkDye ItemID = 1130
	ItemGrayDye ItemID = 1131
	ItemLightGrayDye ItemID = 1132
	ItemCyanDye ItemID = 1133
	ItemPurpleDye ItemID = 1134
	ItemBlueDye ItemID = 1135
	ItemBrownDye ItemID = 1136
	ItemGreenDye ItemID = 1137
	ItemRedDye ItemID = 1138
	ItemBlackDye ItemID = 1139
	ItemBeeSpawnEgg ItemID = 1140
	ItemBlazeSpawnEgg ItemID = 1141
	ItemCatSpawnEgg ItemID = 1142
	ItemCaveSpiderSpawnEgg ItemID = 1143
	ItemChickenSpawnEgg ItemID = 1144
	ItemCodSpawnEgg ItemID = 1145
	ItemCowSpawnEgg ItemID = 1146
	ItemCreeperSpawnEgg ItemID = 1147
	ItemDolphinSpawnEgg ItemID = 1148
	ItemDonkeySpawnEgg ItemID = 1149
	ItemDrownedSpawnEgg ItemID = 1150
	ItemElderGuardianSpawnEgg ItemID = 1151
	ItemEnderDragonSpawnEgg ItemID = 1152
	ItemEndermanSpawnEgg ItemID = 1153
	ItemEndermiteSpawnEgg ItemID = 1154
	ItemEvokerSpawnEgg ItemID = 1155
	ItemFoxSpawnEgg ItemID = 1156
	ItemFrogSpawnEgg ItemID = 1157
	ItemGhastSpawnEgg ItemID = 1158
	ItemGlowSquidSpawnEgg ItemID = 1159
	ItemGoatSpawnEgg ItemID = 1160
	ItemGuardianSpawnEgg ItemID = 1161
	ItemHoglinSpawnEgg ItemID = 1162
	ItemHorseSpawnEgg ItemID = 1163
	ItemHuskSpawnEgg ItemID = 1164
	ItemLlamaSpawnEgg ItemID = 1165
	ItemMagmaCubeSpawnEgg ItemID = 1166
	ItemMooshroomSpawnEgg ItemID = 1167
	ItemOcelotSpawnEgg ItemID = 1168
	ItemPandaSpawnEgg ItemID = 1169
	ItemParrotSpawnEgg ItemID = 1170
	ItemPhantomSpawnEgg ItemID = 1171
	ItemPigSpawnEgg ItemID = 1172
	ItemPiglinSpawnEgg ItemID = 1173
	ItemPiglinBruteSpawnEgg ItemID = 1174
	ItemPillagerSpawnEgg ItemID = 1175
	ItemPolarBearSpawnEgg ItemID = 1176
	ItemPufferfishSpawnEgg ItemID = 1177
	ItemRabbitSpawnEgg ItemID = 1178
	ItemRavagerSpawnEgg ItemID = 1179
	ItemSalmonSpawnEgg ItemID = 1180
	ItemSheepSpawnEgg ItemID = 1181
	ItemShulkerSpawnEgg ItemID = 1182
	ItemSilverfishSpawnEgg ItemID = 1183
	ItemSkeletonSpawnEgg ItemID = 1184
	ItemSkeletonHorseSpawnEgg ItemID = 1185
	ItemSlimeSpawnEgg ItemID = 1186
	ItemSnifferSpawnEgg ItemID = 1187
	ItemSnowGolemSpawnEgg ItemID = 1188
	ItemSpiderSpawnEgg ItemID = 1189
	ItemSquidSpawnEgg ItemID = 1190
	ItemStraySpawnEgg ItemID = 1191
	ItemStriderSpawnEgg ItemID = 1192
	ItemTadpoleSpawnEgg ItemID = 1193
	ItemTraderLlamaSpawnEgg ItemID = 1194
	ItemTropicalFishSpawnEgg ItemID = 1195
	ItemTurtleSpawnEgg ItemID = 1196
	ItemVexSpawnEgg ItemID = 1197
	ItemVillagerSpawnEgg ItemID = 1198
	ItemVindicatorSpawnEgg ItemID = 1199
	ItemWanderingTraderSpawnEgg ItemID = 1200
	ItemWardenSpawnEgg ItemID = 1201
	ItemWitchSpawnEgg ItemID = 1202
	ItemWitherSpawnEgg ItemID = 1203
	ItemWitherSkeletonSpawnEgg ItemID = 1204
	ItemWolfSpawnEgg ItemID = 1205
	ItemZoglinSpawnEgg ItemID = 1206
	ItemZombieSpawnEgg ItemID = 1207
	ItemZombieHorseSpawnEgg ItemID = 1208
	ItemZombieVillagerSpawnEgg ItemID = 1209
	ItemZombifiedPiglinSpawnEgg ItemID = 1210
	ItemExperienceBottle ItemID = 1211
	ItemDebugStick ItemID = 1212
	ItemKnowledgeBook ItemID = 1213
	ItemNetheriteScrap ItemID = 1214
	ItemPrismarineShard ItemID = 1215
	ItemPrismarineCrystals ItemID = 1216
	ItemCopperSulfate ItemID = 1217
	ItemNetheriteItem ItemID = 1218
	ItemHeartOfTheSeaItem ItemID = 1219
	ItemNautilusShellItem ItemID = 1220
	ItemShulkerShellItem ItemID = 1221
	ItemBannerPatternItem ItemID = 1222
	ItemDiscFragment13Item ItemID = 1223
	ItemDiscFragmentItemItem ItemID = 1224
	ItemMusicDisc13Item ItemID = 1225
	ItemMusicDiscCatItem ItemID = 1226
	ItemMusicDiscBlocksItem ItemID = 1227
	ItemMusicDiscChirpItem ItemID = 1228
	ItemMusicDiscFarItem ItemID = 1229
	ItemMusicDiscMallItem ItemID = 1230
	ItemMusicDiscMellohiItem ItemID = 1231
	ItemMusicDiscStalItem ItemID = 1232
	ItemMusicDiscStradItem ItemID = 1233
	ItemMusicDiscWardItem ItemID = 1234
	ItemMusicDisc11Item ItemID = 1235
	ItemMusicDiscWaitItem ItemID = 1236
	ItemMusicDiscOthersideItem ItemID = 1237
	ItemMusicDisc5Item ItemID = 1238
	ItemMusicDiscPigstepItem ItemID = 1239
	ItemMusicDiscRelicItem ItemID = 1240
	ItemMusicDiscDiscFragmentItem ItemID = 1241
	ItemMusicDiscCreatorItem ItemID = 1242
	ItemMusicDiscCreatorMusicBoxItem ItemID = 1243
	ItemGlowInkSacItem ItemID = 1244
	ItemInkSacItem ItemID = 1245
	ItemCocoaBeansItem ItemID = 1246
	ItemWhiteDyeItem ItemID = 1247
	ItemOrangeDyeItem ItemID = 1248
	ItemMagentaDyeItem ItemID = 1249
	ItemLightBlueDyeItem ItemID = 1250
	ItemYellowDyeItem ItemID = 1251
	ItemLimeDyeItem ItemID = 1252
	ItemPinkDyeItem ItemID = 1253
	ItemGrayDyeItem ItemID = 1254
	ItemLightGrayDyeItem ItemID = 1255
	ItemCyanDyeItem ItemID = 1256
	ItemPurpleDyeItem ItemID = 1257
	ItemBlueDyeItem ItemID = 1258
	ItemBrownDyeItem ItemID = 1259
	ItemGreenDyeItem ItemID = 1260
	ItemRedDyeItem ItemID = 1261
	ItemBlackDyeItem ItemID = 1262
)

// initializeItems initializes the item registry.
func initializeItems() {
	itemRegistry = make(map[ItemID]*ItemProperties)
	itemRegistryByName = make(map[string]ItemID)

	// Air
	registerItem(ItemAir, "minecraft:air", &ItemProperties{
		Name:         "Air",
		MaxStackSize: 0,
	})

	// Stone
	registerItem(ItemStone, "minecraft:stone", &ItemProperties{
		Name:         "Stone",
		MaxStackSize: 64,
		BlockID:      BlockStone,
	})

	// Iron Ingot
	registerItem(ItemIronIngot, "minecraft:iron_ingot", &ItemProperties{
		Name:         "Iron Ingot",
		MaxStackSize: 64,
	})

	// Gold Ingot
	registerItem(ItemGoldIngot, "minecraft:gold_ingot", &ItemProperties{
		Name:         "Gold Ingot",
		MaxStackSize: 64,
	})

	// Diamond
	registerItem(ItemDiamond, "minecraft:diamond", &ItemProperties{
		Name:         "Diamond",
		MaxStackSize: 64,
	})

	// Netherite Ingot
	registerItem(ItemNetheriteIngot, "minecraft:netherite_ingot", &ItemProperties{
		Name:         "Netherite Ingot",
		MaxStackSize: 64,
	})

	// Coal
	registerItem(ItemCoal, "minecraft:coal", &ItemProperties{
		Name:         "Coal",
		MaxStackSize: 64,
	})

	// Charcoal
	registerItem(ItemCharcoal, "minecraft:charcoal", &ItemProperties{
		Name:         "Charcoal",
		MaxStackSize: 64,
	})

	// Stick
	registerItem(ItemStick, "minecraft:stick", &ItemProperties{
		Name:         "Stick",
		MaxStackSize: 64,
	})

	// Iron Pickaxe
	registerItem(ItemIronPickaxe, "minecraft:iron_pickaxe", &ItemProperties{
		Name:          "Iron Pickaxe",
		MaxStackSize:  1,
		MaxDurability: 250,
		Tool:          "pickaxe",
		ToolLevel:     2,
		AttackDamage:  2.0,
		AttackSpeed:   1.0,
	})

	// Diamond Sword
	registerItem(ItemDiamondSword, "minecraft:diamond_sword", &ItemProperties{
		Name:          "Diamond Sword",
		MaxStackSize:  1,
		MaxDurability: 1561,
		Tool:          "sword",
		ToolLevel:     3,
		AttackDamage:  7.0,
		AttackSpeed:   1.6,
	})

	// Netherite Sword
	registerItem(ItemNetheriteSword, "minecraft:netherite_sword", &ItemProperties{
		Name:          "Netherite Sword",
		MaxStackSize:  1,
		MaxDurability: 2031,
		Tool:          "sword",
		ToolLevel:     4,
		AttackDamage:  8.0,
		AttackSpeed:   1.6,
	})

	// Bow
	registerItem(ItemBow, "minecraft:bow", &ItemProperties{
		Name:          "Bow",
		MaxStackSize:  1,
		MaxDurability: 384,
		Tool:          "bow",
	})

	// Arrow
	registerItem(ItemArrow, "minecraft:arrow", &ItemProperties{
		Name:         "Arrow",
		MaxStackSize: 64,
	})

	// Apple
	registerItem(ItemApple, "minecraft:apple", &ItemProperties{
		Name:         "Apple",
		MaxStackSize: 64,
		Edible:       true,
		Saturation:   0.3,
	})

	// Golden Apple
	registerItem(ItemGoldenApple, "minecraft:golden_apple", &ItemProperties{
		Name:         "Golden Apple",
		MaxStackSize: 64,
		Edible:       true,
		Saturation:   1.2,
	})

	// Enchanted Golden Apple
	registerItem(ItemEnchantedGoldenApple, "minecraft:enchanted_golden_apple", &ItemProperties{
		Name:         "Enchanted Golden Apple",
		MaxStackSize: 64,
		Edible:       true,
		Saturation:   1.2,
	})

	// Bread
	registerItem(ItemBread, "minecraft:bread", &ItemProperties{
		Name:         "Bread",
		MaxStackSize: 64,
		Edible:       true,
		Saturation:   0.6,
	})

	// Cooked Porkchop
	registerItem(ItemCookedPorkchop, "minecraft:cooked_porkchop", &ItemProperties{
		Name:         "Cooked Porkchop",
		MaxStackSize: 64,
		Edible:       true,
		Saturation:   1.6,
	})

	// Diamond Helmet
	registerItem(ItemDiamondHelmet, "minecraft:diamond_helmet", &ItemProperties{
		Name:          "Diamond Helmet",
		MaxStackSize:  1,
		MaxDurability: 363,
		ArmorMaterial: "diamond",
		Protection:    3.0,
	})

	// Diamond Chestplate
	registerItem(ItemDiamondChestplate, "minecraft:diamond_chestplate", &ItemProperties{
		Name:          "Diamond Chestplate",
		MaxStackSize:  1,
		MaxDurability: 528,
		ArmorMaterial: "diamond",
		Protection:    8.0,
	})

	// Elytra
	registerItem(ItemElytra, "minecraft:elytra", &ItemProperties{
		Name:          "Elytra",
		MaxStackSize:  1,
		MaxDurability: 432,
	})

	// Totem of Undying
	registerItem(ItemTotemOfUndying, "minecraft:totem_of_undying", &ItemProperties{
		Name:         "Totem of Undying",
		MaxStackSize: 1,
	})

	// Shield
	registerItem(ItemShield, "minecraft:shield", &ItemProperties{
		Name:          "Shield",
		MaxStackSize:  1,
		MaxDurability: 336,
	})

	itemCount = len(itemRegistry)
}

// registerItem registers an item in the registry.
func registerItem(id ItemID, name string, props *ItemProperties) {
	props.ID = id
	itemRegistry[id] = props
	itemRegistryByName[name] = id
}

// GetItemProperties returns item properties by ID.
func GetItemProperties(id ItemID) (*ItemProperties, bool) {
	itemMutex.RLock()
	defer itemMutex.RUnlock()

	props, ok := itemRegistry[id]
	return props, ok
}

// GetItemByName returns item ID by name.
func GetItemByName(name string) (ItemID, bool) {
	itemMutex.RLock()
	defer itemMutex.RUnlock()

	id, ok := itemRegistryByName[name]
	return id, ok
}

// GetItemName returns the name of an item by ID.
func GetItemName(id ItemID) (string, bool) {
	itemMutex.RLock()
	defer itemMutex.RUnlock()

	if props, ok := itemRegistry[id]; ok {
		return props.Name, true
	}

	return "", false
}

// ItemCount returns the number of registered items.
func ItemCount() int {
	itemMutex.RLock()
	defer itemMutex.RUnlock()

	return itemCount
}

// GetAllItems returns all registered items.
func GetAllItems() map[ItemID]*ItemProperties {
	itemMutex.RLock()
	defer itemMutex.RUnlock()

	result := make(map[ItemID]*ItemProperties, len(itemRegistry))
	for id, props := range itemRegistry {
		result[id] = props
	}

	return result
}

// GetMaxStackSize returns the maximum stack size for an item.
func GetMaxStackSize(id ItemID) int32 {
	if props, ok := GetItemProperties(id); ok {
		return props.MaxStackSize
	}
	return 64
}

// GetMaxDurability returns the maximum durability for an item.
func GetMaxDurability(id ItemID) int32 {
	if props, ok := GetItemProperties(id); ok {
		return props.MaxDurability
	}
	return 0
}

// IsEdible checks if an item is edible.
func IsEdible(id ItemID) bool {
	if props, ok := GetItemProperties(id); ok {
		return props.Edible
	}
	return false
}

// GetBlockFromItem returns the block ID for an item (if placeable).
func GetBlockFromItem(id ItemID) (BlockID, bool) {
	if props, ok := GetItemProperties(id); ok {
		if props.BlockID != 0 {
			return props.BlockID, true
		}
	}
	return 0, false
}

// NewItem creates a new item stack.
func NewItem(id ItemID, count int32) *Item {
	return &Item{
		ID:      id,
		Count:   count,
		Damage:  0,
		NBTData: nil,
	}
}

// IsEmpty checks if the item stack is empty.
func (i *Item) IsEmpty() bool {
	return i == nil || i.ID == 0 || i.Count <= 0
}

// Clone creates a copy of the item stack.
func (i *Item) Clone() *Item {
	if i == nil {
		return nil
	}

	nbtCopy := make(map[string]interface{})
	for k, v := range i.NBTData {
		nbtCopy[k] = v
	}

	return &Item{
		ID:      i.ID,
		Count:   i.Count,
		Damage:  i.Damage,
		NBTData: nbtCopy,
	}
}

// Split splits the item stack into two stacks.
func (i *Item) Split(amount int32) *Item {
	if i == nil || amount <= 0 || amount > i.Count {
		return nil
	}

	i.Count -= amount
	return &Item{
		ID:      i.ID,
		Count:   amount,
		Damage:  i.Damage,
		NBTData: i.NBTData,
	}
}

// IsDamageable checks if the item can take damage.
func (i *Item) IsDamageable() bool {
	if props, ok := GetItemProperties(i.ID); ok {
		return props.MaxDurability > 0
	}
	return false
}

// IsDamaged checks if the item is damaged.
func (i *Item) IsDamaged() bool {
	return i.IsDamageable() && i.Damage > 0
}

// GetRemainingDurability returns the remaining durability.
func (i *Item) GetRemainingDurability() int32 {
	if !i.IsDamageable() {
		return 0
	}

	if props, ok := GetItemProperties(i.ID); ok {
		return props.MaxDurability - i.Damage
	}

	return 0
}
