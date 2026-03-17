// Package registry provides Minecraft 1.20.1 block and item registries.
package registry

import (
	"sync"
)

// BlockID represents a block state ID.
type BlockID int32

// BlockProperties holds block properties.
type BlockProperties struct {
	ID         BlockID
	Name       string
	Solid      bool
	Transparent bool
	Hardness   float32
	Resistance float32
	Tool       string
	ToolLevel  int
	LightLevel int
	Flammable  bool
	CanBurn    bool
	Replaceable bool
	Slippery   bool
}

// BlockState represents a specific block state.
type BlockState struct {
	ID       BlockID
	Name     string
	Properties map[string]interface{}
}

var (
	blockRegistry      map[BlockID]*BlockProperties
	blockRegistryByName map[string]BlockID
	blockStateRegistry map[BlockID]*BlockState
	blockMutex         sync.RWMutex
	blockCount         int
)

func init() {
	initializeBlocks()
}

// Block constants for Minecraft 1.20.1
const (
	BlockAir BlockID = 0
	BlockStone BlockID = 1
	BlockGranite BlockID = 2
	BlockPolishedGranite BlockID = 3
	BlockDiorite BlockID = 4
	BlockPolishedDiorite BlockID = 5
	BlockAndesite BlockID = 6
	BlockPolishedAndesite BlockID = 7
	BlockDeepslate BlockID = 8
	BlockCobbledDeepslate BlockID = 9
	BlockPolishedDeepslate BlockID = 10
	BlockCalcite BlockID = 11
	BlockTuff BlockID = 12
	BlockDirt BlockID = 13
	BlockCoarseDirt BlockID = 14
	BlockPodzol BlockID = 15
	BlockCrimsonNylium BlockID = 16
	BlockWarpedNylium BlockID = 17
	BlockGrassBlock BlockID = 18
	BlockMycelium BlockID = 19
	BlockRootedDirt BlockID = 20
	BlockMossBlock BlockID = 21
	BlockMossCarpet BlockID = 22
	BlockWater BlockID = 23
	BlockLava BlockID = 24
	BlockSand BlockID = 25
	BlockRedSand BlockID = 26
	BlockGravel BlockID = 27
	BlockSandstone BlockID = 28
	BlockChiseledSandstone BlockID = 29
	BlockCutSandstone BlockID = 30
	BlockSmoothSandstone BlockID = 31
	BlockRedSandstone BlockID = 32
	BlockChiseledRedSandstone BlockID = 33
	BlockCutRedSandstone BlockID = 34
	BlockSmoothRedSandstone BlockID = 35
	BlockClay BlockID = 36
	BlockMud BlockID = 37
	BlockPackedMud BlockID = 38
	BlockMudBricks BlockID = 39
	BlockOakLog BlockID = 40
	BlockSpruceLog BlockID = 41
	BlockBirchLog BlockID = 42
	BlockJungleLog BlockID = 43
	BlockAcaciaLog BlockID = 44
	BlockDarkOakLog BlockID = 45
	BlockCrimsonStem BlockID = 46
	BlockWarpedStem BlockID = 47
	BlockStrippedOakLog BlockID = 48
	BlockStrippedSpruceLog BlockID = 49
	BlockStrippedBirchLog BlockID = 50
	BlockStrippedJungleLog BlockID = 51
	BlockStrippedAcaciaLog BlockID = 52
	BlockStrippedDarkOakLog BlockID = 53
	BlockStrippedCrimsonStem BlockID = 54
	BlockStrippedWarpedStem BlockID = 55
	BlockOakWood BlockID = 56
	BlockSpruceWood BlockID = 57
	BlockBirchWood BlockID = 58
	BlockJungleWood BlockID = 59
	BlockAcaciaWood BlockID = 60
	BlockDarkOakWood BlockID = 61
	BlockCrimsonHyphae BlockID = 62
	BlockWarpedHyphae BlockID = 63
	BlockStrippedOakWood BlockID = 64
	BlockStrippedSpruceWood BlockID = 65
	BlockStrippedBirchWood BlockID = 66
	BlockStrippedJungleWood BlockID = 67
	BlockStrippedAcaciaWood BlockID = 68
	BlockStrippedDarkOakWood BlockID = 69
	BlockStrippedCrimsonHyphae BlockID = 70
	BlockStrippedWarpedHyphae BlockID = 71
	BlockOakPlanks BlockID = 72
	BlockSprucePlanks BlockID = 73
	BlockBirchPlanks BlockID = 74
	BlockJunglePlanks BlockID = 75
	BlockAcaciaPlanks BlockID = 76
	BlockDarkOakPlanks BlockID = 77
	BlockCrimsonPlanks BlockID = 78
	BlockWarpedPlanks BlockID = 79
	BlockBambooPlanks BlockID = 80
	BlockMangrovePlanks BlockID = 81
	BlockCherryPlanks BlockID = 82
	BlockBambooMosaic BlockID = 83
	BlockOakSlab BlockID = 84
	BlockSpruceSlab BlockID = 85
	BlockBirchSlab BlockID = 86
	BlockJungleSlab BlockID = 87
	BlockAcaciaSlab BlockID = 88
	BlockDarkOakSlab BlockID = 89
	BlockCrimsonSlab BlockID = 90
	BlockWarpedSlab BlockID = 91
	BlockBambooSlab BlockID = 92
	BlockMangroveSlab BlockID = 93
	BlockCherrySlab BlockID = 94
	BlockBambooMosaicSlab BlockID = 95
	BlockOakStairs BlockID = 96
	BlockSpruceStairs BlockID = 97
	BlockBirchStairs BlockID = 98
	BlockJungleStairs BlockID = 99
	BlockAcaciaStairs BlockID = 100
	BlockDarkOakStairs BlockID = 101
	BlockCrimsonStairs BlockID = 102
	BlockWarpedStairs BlockID = 103
	BlockBambooStairs BlockID = 104
	BlockMangroveStairs BlockID = 105
	BlockCherryStairs BlockID = 106
	BlockGlass BlockID = 107
	BlockTintedGlass BlockID = 108
	BlockWhiteStainedGlass BlockID = 109
	BlockOrangeStainedGlass BlockID = 110
	BlockMagentaStainedGlass BlockID = 111
	BlockLightBlueStainedGlass BlockID = 112
	BlockYellowStainedGlass BlockID = 113
	BlockLimeStainedGlass BlockID = 114
	BlockPinkStainedGlass BlockID = 115
	BlockGrayStainedGlass BlockID = 116
	BlockLightGrayStainedGlass BlockID = 117
	BlockCyanStainedGlass BlockID = 118
	BlockPurpleStainedGlass BlockID = 119
	BlockBlueStainedGlass BlockID = 120
	BlockBrownStainedGlass BlockID = 121
	BlockGreenStainedGlass BlockID = 122
	BlockRedStainedGlass BlockID = 123
	BlockBlackStainedGlass BlockID = 124
	BlockGlassPane BlockID = 125
	BlockWhiteStainedGlassPane BlockID = 126
	BlockOrangeStainedGlassPane BlockID = 127
	BlockMagentaStainedGlassPane BlockID = 128
	BlockLightBlueStainedGlassPane BlockID = 129
	BlockYellowStainedGlassPane BlockID = 130
	BlockLimeStainedGlassPane BlockID = 131
	BlockPinkStainedGlassPane BlockID = 132
	BlockGrayStainedGlassPane BlockID = 133
	BlockLightGrayStainedGlassPane BlockID = 134
	BlockCyanStainedGlassPane BlockID = 135
	BlockPurpleStainedGlassPane BlockID = 136
	BlockBlueStainedGlassPane BlockID = 137
	BlockBrownStainedGlassPane BlockID = 138
	BlockGreenStainedGlassPane BlockID = 139
	BlockRedStainedGlassPane BlockID = 140
	BlockBlackStainedGlassPane BlockID = 141
	BlockPrismarine BlockID = 142
	BlockPrismarineBricks BlockID = 143
	BlockDarkPrismarine BlockID = 144
	BlockSeaLantern BlockID = 145
	BlockCopperBlock BlockID = 146
	BlockExposedCopper BlockID = 147
	BlockWeatheredCopper BlockID = 148
	BlockOxidizedCopper BlockID = 149
	BlockCutCopper BlockID = 150
	BlockExposedCutCopper BlockID = 151
	BlockWeatheredCutCopper BlockID = 152
	BlockOxidizedCutCopper BlockID = 153
	BlockChiseledCopper BlockID = 154
	BlockExposedChiseledCopper BlockID = 155
	BlockWeatheredChiseledCopper BlockID = 156
	BlockOxidizedChiseledCopper BlockID = 157
	BlockCopperGrate BlockID = 158
	BlockExposedCopperGrate BlockID = 159
	BlockWeatheredCopperGrate BlockID = 160
	BlockOxidizedCopperGrate BlockID = 161
	BlockCopperBulb BlockID = 162
	BlockExposedCopperBulb BlockID = 163
	BlockWeatheredCopperBulb BlockID = 164
	BlockOxidizedCopperBulb BlockID = 165
	BlockCopperDoor BlockID = 166
	BlockExposedCopperDoor BlockID = 167
	BlockWeatheredCopperDoor BlockID = 168
	BlockOxidizedCopperDoor BlockID = 169
	BlockCopperTrapdoor BlockID = 170
	BlockExposedCopperTrapdoor BlockID = 171
	BlockWeatheredCopperTrapdoor BlockID = 172
	BlockOxidizedCopperTrapdoor BlockID = 173
	BlockCopperGrateSlab BlockID = 174
	BlockExposedCopperGrateSlab BlockID = 175
	BlockWeatheredCopperGrateSlab BlockID = 176
	BlockOxidizedCopperGrateSlab BlockID = 177
	BlockWaxedCopper BlockID = 178
	BlockWaxedExposedCopper BlockID = 179
	BlockWaxedWeatheredCopper BlockID = 180
	BlockWaxedOxidizedCopper BlockID = 181
	BlockWaxedCutCopper BlockID = 182
	BlockWaxedExposedCutCopper BlockID = 183
	BlockWaxedWeatheredCutCopper BlockID = 184
	BlockWaxedOxidizedCutCopper BlockID = 185
	BlockWaxedChiseledCopper BlockID = 186
	BlockWaxedExposedChiseledCopper BlockID = 187
	BlockWaxedWeatheredChiseledCopper BlockID = 188
	BlockWaxedOxidizedChiseledCopper BlockID = 189
	BlockWaxedCopperGrate BlockID = 190
	BlockWaxedExposedCopperGrate BlockID = 191
	BlockWaxedWeatheredCopperGrate BlockID = 192
	BlockWaxedOxidizedCopperGrate BlockID = 193
	BlockWaxedCopperBulb BlockID = 194
	BlockWaxedExposedCopperBulb BlockID = 195
	BlockWaxedWeatheredCopperBulb BlockID = 196
	BlockWaxedOxidizedCopperBulb BlockID = 197
	BlockWaxedCopperDoor BlockID = 198
	BlockWaxedExposedCopperDoor BlockID = 199
	BlockWaxedWeatheredCopperDoor BlockID = 200
	BlockWaxedOxidizedCopperDoor BlockID = 201
	BlockWaxedCopperTrapdoor BlockID = 202
	BlockWaxedExposedCopperTrapdoor BlockID = 203
	BlockWaxedWeatheredCopperTrapdoor BlockID = 204
	BlockWaxedOxidizedCopperTrapdoor BlockID = 205
	BlockWaxedCopperGrateSlab BlockID = 206
	BlockWaxedExposedCopperGrateSlab BlockID = 207
	BlockWaxedWeatheredCopperGrateSlab BlockID = 208
	BlockWaxedOxidizedCopperGrateSlab BlockID = 209
	BlockIronBlock BlockID = 210
	BlockCauldron BlockID = 211
	BlockWaterCauldron BlockID = 212
	BlockLavaCauldron BlockID = 213
	BlockPowderSnowCauldron BlockID = 214
	BlockCopperSulfateCauldron BlockID = 215
	BlockGoldBlock BlockID = 216
	BlockRawIronBlock BlockID = 217
	BlockRawCopperBlock BlockID = 218
	BlockRawGoldBlock BlockID = 219
	BlockAmethystBlock BlockID = 220
	BlockBuddingAmethyst BlockID = 221
	BlockCopperOre BlockID = 222
	BlockDeepslateCopperOre BlockID = 223
	BlockIronOre BlockID = 224
	BlockDeepslateIronOre BlockID = 225
	BlockGoldOre BlockID = 226
	BlockDeepslateGoldOre BlockID = 227
	BlockDiamondOre BlockID = 228
	BlockDeepslateDiamondOre BlockID = 229
	BlockLapisOre BlockID = 230
	BlockDeepslateLapisOre BlockID = 231
	BlockRedstoneOre BlockID = 232
	BlockDeepslateRedstoneOre BlockID = 233
	BlockEmeraldOre BlockID = 234
	BlockDeepslateEmeraldOre BlockID = 235
	BlockCoalOre BlockID = 236
	BlockDeepslateCoalOre BlockID = 237
	BlockNetherQuartzOre BlockID = 238
	BlockNetherGoldOre BlockID = 239
	BlockAncientDebris BlockID = 240
	BlockCoalBlock BlockID = 241
	BlockDiamondBlock BlockID = 242
	BlockLapisBlock BlockID = 243
	BlockEmeraldBlock BlockID = 244
	BlockRedstoneBlock BlockID = 245
	BlockNetheriteBlock BlockID = 246
	BlockOakFence BlockID = 247
	BlockSpruceFence BlockID = 248
	BlockBirchFence BlockID = 249
	BlockJungleFence BlockID = 250
	BlockAcaciaFence BlockID = 251
	BlockDarkOakFence BlockID = 252
	BlockCrimsonFence BlockID = 253
	BlockWarpedFence BlockID = 254
	BlockMangroveFence BlockID = 255
	BlockCherryFence BlockID = 256
	BlockBambooFence BlockID = 257
	BlockOakFenceGate BlockID = 258
	BlockSpruceFenceGate BlockID = 259
	BlockBirchFenceGate BlockID = 260
	BlockJungleFenceGate BlockID = 261
	BlockAcaciaFenceGate BlockID = 262
	BlockDarkOakFenceGate BlockID = 263
	BlockCrimsonFenceGate BlockID = 264
	BlockWarpedFenceGate BlockID = 265
	BlockMangroveFenceGate BlockID = 266
	BlockCherryFenceGate BlockID = 267
	BlockBambooFenceGate BlockID = 268
	BlockOakDoor BlockID = 269
	BlockSpruceDoor BlockID = 270
	BlockBirchDoor BlockID = 271
	BlockJungleDoor BlockID = 272
	BlockAcaciaDoor BlockID = 273
	BlockDarkOakDoor BlockID = 274
	BlockCrimsonDoor BlockID = 275
	BlockWarpedDoor BlockID = 276
	BlockMangroveDoor BlockID = 277
	BlockCherryDoor BlockID = 278
	BlockBambooDoor BlockID = 279
	BlockIronDoor BlockID = 280
	BlockCopperDoorBlock BlockID = 281
	BlockExposedCopperDoorBlock BlockID = 282
	BlockWeatheredCopperDoorBlock BlockID = 283
	BlockOxidizedCopperDoorBlock BlockID = 284
	BlockWaxedCopperDoorBlock BlockID = 285
	BlockWaxedExposedCopperDoorBlock BlockID = 286
	BlockWaxedWeatheredCopperDoorBlock BlockID = 287
	BlockWaxedOxidizedCopperDoorBlock BlockID = 288
	BlockOakTrapdoor BlockID = 289
	BlockSpruceTrapdoor BlockID = 290
	BlockBirchTrapdoor BlockID = 291
	BlockJungleTrapdoor BlockID = 292
	BlockAcaciaTrapdoor BlockID = 293
	BlockDarkOakTrapdoor BlockID = 294
	BlockCrimsonTrapdoor BlockID = 295
	BlockWarpedTrapdoor BlockID = 296
	BlockMangroveTrapdoor BlockID = 297
	BlockCherryTrapdoor BlockID = 298
	BlockBambooTrapdoor BlockID = 299
	BlockIronTrapdoor BlockID = 300
	BlockCopperTrapdoorBlock BlockID = 301
	BlockExposedCopperTrapdoorBlock BlockID = 302
	BlockWeatheredCopperTrapdoorBlock BlockID = 303
	BlockOxidizedCopperTrapdoorBlock BlockID = 304
	BlockWaxedCopperTrapdoorBlock BlockID = 305
	BlockWaxedExposedCopperTrapdoorBlock BlockID = 306
	BlockWaxedWeatheredCopperTrapdoorBlock BlockID = 307
	BlockWaxedOxidizedCopperTrapdoorBlock BlockID = 308
	BlockOakPressurePlate BlockID = 309
	BlockSprucePressurePlate BlockID = 310
	BlockBirchPressurePlate BlockID = 311
	BlockJunglePressurePlate BlockID = 312
	BlockAcaciaPressurePlate BlockID = 313
	BlockDarkOakPressurePlate BlockID = 314
	BlockCrimsonPressurePlate BlockID = 315
	BlockWarpedPressurePlate BlockID = 316
	BlockMangrovePressurePlate BlockID = 317
	BlockCherryPressurePlate BlockID = 318
	BlockBambooPressurePlate BlockID = 319
	BlockStonePressurePlate BlockID = 320
	BlockPolishedBlackstonePressurePlate BlockID = 321
	BlockHeavyWeightedPressurePlate BlockID = 322
	BlockLightWeightedPressurePlate BlockID = 323
	BlockCopperPressurePlate BlockID = 324
	BlockExposedCopperPressurePlate BlockID = 325
	BlockWeatheredCopperPressurePlate BlockID = 326
	BlockOxidizedCopperPressurePlate BlockID = 327
	BlockWaxedCopperPressurePlate BlockID = 328
	BlockWaxedExposedCopperPressurePlate BlockID = 329
	BlockWaxedWeatheredCopperPressurePlate BlockID = 330
	BlockWaxedOxidizedCopperPressurePlate BlockID = 331
	BlockOakButton BlockID = 332
	BlockSpruceButton BlockID = 333
	BlockBirchButton BlockID = 334
	BlockJungleButton BlockID = 335
	BlockAcaciaButton BlockID = 336
	BlockDarkOakButton BlockID = 337
	BlockCrimsonButton BlockID = 338
	BlockWarpedButton BlockID = 339
	BlockMangroveButton BlockID = 340
	BlockCherryButton BlockID = 341
	BlockBambooButton BlockID = 342
	BlockStoneButton BlockID = 343
	BlockPolishedBlackstoneButton BlockID = 344
	BlockCopperButton BlockID = 345
	BlockExposedCopperButton BlockID = 346
	BlockWeatheredCopperButton BlockID = 347
	BlockOxidizedCopperButton BlockID = 348
	BlockWaxedCopperButton BlockID = 349
	BlockWaxedExposedCopperButton BlockID = 350
	BlockWaxedWeatheredCopperButton BlockID = 351
	BlockWaxedOxidizedCopperButton BlockID = 352
	BlockBricks BlockID = 353
	BlockBrickSlab BlockID = 354
	BlockBrickStairs BlockID = 355
	BlockBrickWall BlockID = 356
	BlockMudBrickSlab BlockID = 357
	BlockMudBrickStairs BlockID = 358
	BlockMudBrickWall BlockID = 359
	BlockCobbledDeepslateSlab BlockID = 360
	BlockCobbledDeepslateStairs BlockID = 361
	BlockCobbledDeepslateWall BlockID = 362
	BlockPolishedDeepslateSlab BlockID = 363
	BlockPolishedDeepslateStairs BlockID = 364
	BlockPolishedDeepslateWall BlockID = 365
	BlockDeepslateBricks BlockID = 366
	BlockDeepslateBrickSlab BlockID = 367
	BlockDeepslateBrickStairs BlockID = 368
	BlockDeepslateBrickWall BlockID = 369
	BlockDeepslateTiles BlockID = 370
	BlockDeepslateTileSlab BlockID = 371
	BlockDeepslateTileStairs BlockID = 372
	BlockDeepslateTileWall BlockID = 373
	BlockReinforcedDeepslate BlockID = 374
	BlockCobblestone BlockID = 375
	BlockCobblestoneSlab BlockID = 376
	BlockCobblestoneStairs BlockID = 377
	BlockCobblestoneWall BlockID = 378
	BlockMossyCobblestone BlockID = 379
	BlockMossyCobblestoneSlab BlockID = 380
	BlockMossyCobblestoneStairs BlockID = 381
	BlockMossyCobblestoneWall BlockID = 382
	BlockStoneBricks BlockID = 383
	BlockStoneBrickSlab BlockID = 384
	BlockStoneBrickStairs BlockID = 385
	BlockStoneBrickWall BlockID = 386
	BlockMossyStoneBricks BlockID = 387
	BlockCrackedStoneBricks BlockID = 388
	BlockChiseledStoneBricks BlockID = 389
	BlockSmoothStone BlockID = 390
	BlockSmoothStoneSlab BlockID = 391
	BlockSandstoneSlab BlockID = 392
	BlockSandstoneStairs BlockID = 393
	BlockSandstoneWall BlockID = 394
	BlockRedSandstoneSlab BlockID = 395
	BlockRedSandstoneStairs BlockID = 396
	BlockRedSandstoneWall BlockID = 397
	BlockPrismarineSlab BlockID = 398
	BlockPrismarineStairs BlockID = 399
	BlockPrismarineWall BlockID = 400
	BlockPrismarineBrickSlab BlockID = 401
	BlockPrismarineBrickStairs BlockID = 402
	BlockDarkPrismarineSlab BlockID = 403
	BlockDarkPrismarineStairs BlockID = 404
	BlockNetherBricks BlockID = 405
	BlockNetherBrickSlab BlockID = 406
	BlockNetherBrickStairs BlockID = 407
	BlockNetherBrickFence BlockID = 408
	BlockNetherBrickWall BlockID = 409
	BlockRedNetherBricks BlockID = 410
	BlockBlackstone BlockID = 411
	BlockBlackstoneSlab BlockID = 412
	BlockBlackstoneStairs BlockID = 413
	BlockBlackstoneWall BlockID = 414
	BlockPolishedBlackstone BlockID = 415
	BlockPolishedBlackstoneSlab BlockID = 416
	BlockPolishedBlackstoneStairs BlockID = 417
	BlockPolishedBlackstoneWall BlockID = 418
	BlockPolishedBlackstoneBricks BlockID = 419
	BlockPolishedBlackstoneBrickSlab BlockID = 420
	BlockPolishedBlackstoneBrickStairs BlockID = 421
	BlockPolishedBlackstoneBrickWall BlockID = 422
	BlockChiseledPolishedBlackstone BlockID = 423
	BlockCrackedPolishedBlackstoneBricks BlockID = 424
	BlockGildedBlackstone BlockID = 425
	BlockEndStoneBricks BlockID = 426
	BlockEndStoneBrickSlab BlockID = 427
	BlockEndStoneBrickStairs BlockID = 428
	BlockEndStoneBrickWall BlockID = 429
	BlockPurpurBlock BlockID = 430
	BlockPurpurPillar BlockID = 431
	BlockPurpurSlab BlockID = 432
	BlockPurpurStairs BlockID = 433
	BlockQuartzBlock BlockID = 434
	BlockChiseledQuartzBlock BlockID = 435
	BlockQuartzPillar BlockID = 436
	BlockQuartzBricks BlockID = 437
	BlockQuartzSlab BlockID = 438
	BlockQuartzStairs BlockID = 439
	BlockSmoothQuartz BlockID = 440
	BlockTerracotta BlockID = 441
	BlockWhiteTerracotta BlockID = 442
	BlockOrangeTerracotta BlockID = 443
	BlockMagentaTerracotta BlockID = 444
	BlockLightBlueTerracotta BlockID = 445
	BlockYellowTerracotta BlockID = 446
	BlockLimeTerracotta BlockID = 447
	BlockPinkTerracotta BlockID = 448
	BlockGrayTerracotta BlockID = 449
	BlockLightGrayTerracotta BlockID = 450
	BlockCyanTerracotta BlockID = 451
	BlockPurpleTerracotta BlockID = 452
	BlockBlueTerracotta BlockID = 453
	BlockBrownTerracotta BlockID = 454
	BlockGreenTerracotta BlockID = 455
	BlockRedTerracotta BlockID = 456
	BlockBlackTerracotta BlockID = 457
	BlockWhiteGlazedTerracotta BlockID = 458
	BlockOrangeGlazedTerracotta BlockID = 459
	BlockMagentaGlazedTerracotta BlockID = 460
	BlockLightBlueGlazedTerracotta BlockID = 461
	BlockYellowGlazedTerracotta BlockID = 462
	BlockLimeGlazedTerracotta BlockID = 463
	BlockPinkGlazedTerracotta BlockID = 464
	BlockGrayGlazedTerracotta BlockID = 465
	BlockLightGrayGlazedTerracotta BlockID = 466
	BlockCyanGlazedTerracotta BlockID = 467
	BlockPurpleGlazedTerracotta BlockID = 468
	BlockBlueGlazedTerracotta BlockID = 469
	BlockBrownGlazedTerracotta BlockID = 470
	BlockGreenGlazedTerracotta BlockID = 471
	BlockRedGlazedTerracotta BlockID = 472
	BlockBlackGlazedTerracotta BlockID = 473
	BlockWhiteWool BlockID = 474
	BlockOrangeWool BlockID = 475
	BlockMagentaWool BlockID = 476
	BlockLightBlueWool BlockID = 477
	BlockYellowWool BlockID = 478
	BlockLimeWool BlockID = 479
	BlockPinkWool BlockID = 480
	BlockGrayWool BlockID = 481
	BlockLightGrayWool BlockID = 482
	BlockCyanWool BlockID = 483
	BlockPurpleWool BlockID = 484
	BlockBlueWool BlockID = 485
	BlockBrownWool BlockID = 486
	BlockGreenWool BlockID = 487
	BlockRedWool BlockID = 488
	BlockBlackWool BlockID = 489
	BlockWhiteCarpet BlockID = 490
	BlockOrangeCarpet BlockID = 491
	BlockMagentaCarpet BlockID = 492
	BlockLightBlueCarpet BlockID = 493
	BlockYellowCarpet BlockID = 494
	BlockLimeCarpet BlockID = 495
	BlockPinkCarpet BlockID = 496
	BlockGrayCarpet BlockID = 497
	BlockLightGrayCarpet BlockID = 498
	BlockCyanCarpet BlockID = 499
	BlockPurpleCarpet BlockID = 500
	BlockBlueCarpet BlockID = 501
	BlockBrownCarpet BlockID = 502
	BlockGreenCarpet BlockID = 503
	BlockRedCarpet BlockID = 504
	BlockBlackCarpet BlockID = 505
	BlockCraftingTable BlockID = 506
	BlockFurnace BlockID = 507
	BlockSmoker BlockID = 508
	BlockBlastFurnace BlockID = 509
	BlockLectern BlockID = 510
	BlockLoom BlockID = 511
	BlockCartographyTable BlockID = 512
	BlockFletchingTable BlockID = 513
	BlockGrindstone BlockID = 514
	BlockSmithingTable BlockID = 515
	BlockStonecutter BlockID = 516
	BlockBrewingStand BlockID = 517
	BlockBarrel BlockID = 518
	BlockChiseledBookshelf BlockID = 519
	BlockDecoratedPot BlockID = 520
	BlockShulkerBox BlockID = 521
	BlockWhiteShulkerBox BlockID = 522
	BlockOrangeShulkerBox BlockID = 523
	BlockMagentaShulkerBox BlockID = 524
	BlockLightBlueShulkerBox BlockID = 525
	BlockYellowShulkerBox BlockID = 526
	BlockLimeShulkerBox BlockID = 527
	BlockPinkShulkerBox BlockID = 528
	BlockGrayShulkerBox BlockID = 529
	BlockLightGrayShulkerBox BlockID = 530
	BlockCyanShulkerBox BlockID = 531
	BlockPurpleShulkerBox BlockID = 532
	BlockBlueShulkerBox BlockID = 533
	BlockBrownShulkerBox BlockID = 534
	BlockGreenShulkerBox BlockID = 535
	BlockRedShulkerBox BlockID = 536
	BlockBlackShulkerBox BlockID = 537
	BlockChest BlockID = 538
	BlockTrappedChest BlockID = 539
	BlockEnderChest BlockID = 540
	BlockBarrelBlock BlockID = 541
	BlockOakSign BlockID = 542
	BlockSpruceSign BlockID = 543
	BlockBirchSign BlockID = 544
	BlockJungleSign BlockID = 545
	BlockAcaciaSign BlockID = 546
	BlockDarkOakSign BlockID = 547
	BlockCrimsonSign BlockID = 548
	BlockWarpedSign BlockID = 549
	BlockMangroveSign BlockID = 550
	BlockCherrySign BlockID = 551
	BlockBambooSign BlockID = 552
	BlockOakHangingSign BlockID = 553
	BlockSpruceHangingSign BlockID = 554
	BlockBirchHangingSign BlockID = 555
	BlockJungleHangingSign BlockID = 556
	BlockAcaciaHangingSign BlockID = 557
	BlockDarkOakHangingSign BlockID = 558
	BlockCrimsonHangingSign BlockID = 559
	BlockWarpedHangingSign BlockID = 560
	BlockMangroveHangingSign BlockID = 561
	BlockCherryHangingSign BlockID = 562
	BlockBambooHangingSign BlockID = 563
	BlockBedrock BlockID = 564
	BlockSandstoneBlock BlockID = 565
	BlockWhiteConcrete BlockID = 566
	BlockOrangeConcrete BlockID = 567
	BlockMagentaConcrete BlockID = 568
	BlockLightBlueConcrete BlockID = 569
	BlockYellowConcrete BlockID = 570
	BlockLimeConcrete BlockID = 571
	BlockPinkConcrete BlockID = 572
	BlockGrayConcrete BlockID = 573
	BlockLightGrayConcrete BlockID = 574
	BlockCyanConcrete BlockID = 575
	BlockPurpleConcrete BlockID = 576
	BlockBlueConcrete BlockID = 577
	BlockBrownConcrete BlockID = 578
	BlockGreenConcrete BlockID = 579
	BlockRedConcrete BlockID = 580
	BlockBlackConcrete BlockID = 581
	BlockWhiteConcretePowder BlockID = 582
	BlockOrangeConcretePowder BlockID = 583
	BlockMagentaConcretePowder BlockID = 584
	BlockLightBlueConcretePowder BlockID = 585
	BlockYellowConcretePowder BlockID = 586
	BlockLimeConcretePowder BlockID = 587
	BlockPinkConcretePowder BlockID = 588
	BlockGrayConcretePowder BlockID = 589
	BlockLightGrayConcretePowder BlockID = 590
	BlockCyanConcretePowder BlockID = 591
	BlockPurpleConcretePowder BlockID = 592
	BlockBlueConcretePowder BlockID = 593
	BlockBrownConcretePowder BlockID = 594
	BlockGreenConcretePowder BlockID = 595
	BlockRedConcretePowder BlockID = 596
	BlockBlackConcretePowder BlockID = 597
	BlockSponge BlockID = 598
	BlockWetSponge BlockID = 599
	BlockGlowstone BlockID = 600
	BlockSeaPickle BlockID = 601
	BlockSnowBlock BlockID = 602
	BlockSnow BlockID = 603
	BlockIce BlockID = 604
	BlockPackedIce BlockID = 605
	BlockBlueIce BlockID = 606
	BlockPowderSnow BlockID = 607
	BlockFire BlockID = 608
	BlockSoulFire BlockID = 609
	BlockCampfire BlockID = 610
	BlockSoulCampfire BlockID = 611
	BlockTorch BlockID = 612
	BlockSoulTorch BlockID = 613
	BlockEndRod BlockID = 614
	BlockLantern BlockID = 615
	BlockSoulLantern BlockID = 616
	BlockGlowLichen BlockID = 617
	BlockLight BlockID = 618
	BlockLightningRod BlockID = 619
	BlockCopperBulbBlock BlockID = 620
	BlockExposedCopperBulbBlock BlockID = 621
	BlockWeatheredCopperBulbBlock BlockID = 622
	BlockOxidizedCopperBulbBlock BlockID = 623
	BlockWaxedCopperBulbBlock BlockID = 624
	BlockWaxedExposedCopperBulbBlock BlockID = 625
	BlockWaxedWeatheredCopperBulbBlock BlockID = 626
	BlockWaxedOxidizedCopperBulbBlock BlockID = 627
	BlockSeaPickleBlock BlockID = 628
	BlockSculk BlockID = 629
	BlockSculkSensor BlockID = 630
	BlockSculkCatalyst BlockID = 631
	BlockSculkVein BlockID = 632
	BlockSculkShrieker BlockID = 633
	BlockNetherPortal BlockID = 634
	BlockEndPortal BlockID = 635
	BlockEndPortalFrame BlockID = 636
	BlockEndGateway BlockID = 637
	BlockCommandBlock BlockID = 638
	BlockChainCommandBlock BlockID = 639
	BlockRepeatingCommandBlock BlockID = 640
	BlockStructureBlock BlockID = 641
	BlockJigsaw BlockID = 642
	BlockSpawner BlockID = 643
	BlockObsidian BlockID = 644
	BlockCryingObsidian BlockID = 645
	BlockDragonEgg BlockID = 646
	BlockNetherStar BlockID = 647
	BlockEnchantingTable BlockID = 648
	BlockBookshelf BlockID = 649
	BlockBeehive BlockID = 650
	BlockNetherrack BlockID = 651
	BlockEndStone BlockID = 652
)

// initializeBlocks initializes the block registry.
func initializeBlocks() {
	blockRegistry = make(map[BlockID]*BlockProperties)
	blockRegistryByName = make(map[string]BlockID)
	blockStateRegistry = make(map[BlockID]*BlockState)

	// Air
	registerBlock(BlockAir, "minecraft:air", &BlockProperties{
		Name:        "Air",
		Transparent: true,
		Replaceable: true,
	})

	// Stone
	registerBlock(BlockStone, "minecraft:stone", &BlockProperties{
		Name:        "Stone",
		Solid:       true,
		Hardness:    1.5,
		Resistance:  6.0,
		Tool:        "pickaxe",
		ToolLevel:   0,
	})

	// Granite
	registerBlock(BlockGranite, "minecraft:granite", &BlockProperties{
		Name:        "Granite",
		Solid:       true,
		Hardness:    1.5,
		Resistance:  6.0,
		Tool:        "pickaxe",
		ToolLevel:   0,
	})

	// Dirt
	registerBlock(BlockDirt, "minecraft:dirt", &BlockProperties{
		Name:        "Dirt",
		Solid:       true,
		Hardness:    0.5,
		Resistance:  0.5,
		Tool:        "shovel",
		ToolLevel:   0,
	})

	// Grass Block
	registerBlock(BlockGrassBlock, "minecraft:grass_block", &BlockProperties{
		Name:        "Grass Block",
		Solid:       true,
		Hardness:    0.6,
		Resistance:  0.6,
		Tool:        "shovel",
		ToolLevel:   0,
	})

	// Water
	registerBlock(BlockWater, "minecraft:water", &BlockProperties{
		Name:        "Water",
		Solid:       false,
		Transparent: true,
		Hardness:    100.0,
		Replaceable: true,
	})

	// Lava
	registerBlock(BlockLava, "minecraft:lava", &BlockProperties{
		Name:        "Lava",
		Solid:       false,
		Transparent: true,
		Hardness:    100.0,
		LightLevel:  15,
	})

	// Sand
	registerBlock(BlockSand, "minecraft:sand", &BlockProperties{
		Name:        "Sand",
		Solid:       true,
		Hardness:    0.5,
		Resistance:  0.5,
		Tool:        "shovel",
		ToolLevel:   0,
	})

	// Gravel
	registerBlock(BlockGravel, "minecraft:gravel", &BlockProperties{
		Name:        "Gravel",
		Solid:       true,
		Hardness:    0.6,
		Resistance:  0.6,
		Tool:        "shovel",
		ToolLevel:   0,
	})

	// Oak Log
	registerBlock(BlockOakLog, "minecraft:oak_log", &BlockProperties{
		Name:        "Oak Log",
		Solid:       true,
		Hardness:    2.0,
		Resistance:  2.0,
		Tool:        "axe",
		ToolLevel:   0,
		Flammable:   true,
		CanBurn:     true,
	})

	// Oak Planks
	registerBlock(BlockOakPlanks, "minecraft:oak_planks", &BlockProperties{
		Name:        "Oak Planks",
		Solid:       true,
		Hardness:    2.0,
		Resistance:  3.0,
		Tool:        "axe",
		ToolLevel:   0,
		Flammable:   true,
		CanBurn:     true,
	})

	// Glass
	registerBlock(BlockGlass, "minecraft:glass", &BlockProperties{
		Name:        "Glass",
		Solid:       true,
		Transparent: true,
		Hardness:    0.3,
		Resistance:  0.3,
	})

	// Crafting Table
	registerBlock(BlockCraftingTable, "minecraft:crafting_table", &BlockProperties{
		Name:        "Crafting Table",
		Solid:       true,
		Hardness:    2.5,
		Resistance:  2.5,
		Tool:        "axe",
		ToolLevel:   0,
		Flammable:   true,
		CanBurn:     true,
	})

	// Furnace
	registerBlock(BlockFurnace, "minecraft:furnace", &BlockProperties{
		Name:        "Furnace",
		Solid:       true,
		Hardness:    3.5,
		Resistance:  3.5,
		Tool:        "pickaxe",
		ToolLevel:   0,
	})

	// Chest
	registerBlock(BlockChest, "minecraft:chest", &BlockProperties{
		Name:        "Chest",
		Solid:       true,
		Hardness:    2.5,
		Resistance:  2.5,
		Tool:        "axe",
		ToolLevel:   0,
		Flammable:   true,
		CanBurn:     true,
	})

	// Bedrock
	registerBlock(BlockBedrock, "minecraft:bedrock", &BlockProperties{
		Name:        "Bedrock",
		Solid:       true,
		Hardness:    -1.0,
		Resistance:  3600000.0,
	})

	// Iron Block
	registerBlock(BlockIronBlock, "minecraft:iron_block", &BlockProperties{
		Name:        "Iron Block",
		Solid:       true,
		Hardness:    5.0,
		Resistance:  6.0,
		Tool:        "pickaxe",
		ToolLevel:   1,
	})

	// Gold Block
	registerBlock(BlockGoldBlock, "minecraft:gold_block", &BlockProperties{
		Name:        "Gold Block",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  6.0,
		Tool:        "pickaxe",
		ToolLevel:   2,
	})

	// Diamond Block
	registerBlock(BlockDiamondBlock, "minecraft:diamond_block", &BlockProperties{
		Name:        "Diamond Block",
		Solid:       true,
		Hardness:    5.0,
		Resistance:  6.0,
		Tool:        "pickaxe",
		ToolLevel:   3,
	})

	// Coal Ore
	registerBlock(BlockCoalOre, "minecraft:coal_ore", &BlockProperties{
		Name:        "Coal Ore",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "pickaxe",
		ToolLevel:   0,
	})

	// Iron Ore
	registerBlock(BlockIronOre, "minecraft:iron_ore", &BlockProperties{
		Name:        "Iron Ore",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "pickaxe",
		ToolLevel:   1,
	})

	// Gold Ore
	registerBlock(BlockGoldOre, "minecraft:gold_ore", &BlockProperties{
		Name:        "Gold Ore",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "pickaxe",
		ToolLevel:   2,
	})

	// Diamond Ore
	registerBlock(BlockDiamondOre, "minecraft:diamond_ore", &BlockProperties{
		Name:        "Diamond Ore",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "pickaxe",
		ToolLevel:   3,
	})

	// Copper Ore
	registerBlock(BlockCopperOre, "minecraft:copper_ore", &BlockProperties{
		Name:        "Copper Ore",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "pickaxe",
		ToolLevel:   1,
	})

	// Deepslate
	registerBlock(BlockDeepslate, "minecraft:deepslate", &BlockProperties{
		Name:        "Deepslate",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "pickaxe",
		ToolLevel:   0,
	})

	// Obsidian
	registerBlock(BlockObsidian, "minecraft:obsidian", &BlockProperties{
		Name:        "Obsidian",
		Solid:       true,
		Hardness:    50.0,
		Resistance:  1200.0,
		Tool:        "pickaxe",
		ToolLevel:   3,
	})

	// Glowstone
	registerBlock(BlockGlowstone, "minecraft:glowstone", &BlockProperties{
		Name:        "Glowstone",
		Solid:       true,
		Transparent: true,
		Hardness:    0.3,
		Resistance:  0.3,
		LightLevel:  15,
	})

	// Sea Lantern
	registerBlock(BlockSeaLantern, "minecraft:sea_lantern", &BlockProperties{
		Name:        "Sea Lantern",
		Solid:       true,
		Transparent: true,
		Hardness:    0.3,
		Resistance:  0.3,
		LightLevel:  15,
	})

	// Shulker Box
	registerBlock(BlockShulkerBox, "minecraft:shulker_box", &BlockProperties{
		Name:        "Shulker Box",
		Solid:       true,
		Hardness:    2.0,
		Resistance:  2.0,
		Tool:        "axe",
		ToolLevel:   0,
	})

	// Beehive
	registerBlock(BlockBeehive, "minecraft:beehive", &BlockProperties{
		Name:        "Beehive",
		Solid:       true,
		Hardness:    0.6,
		Resistance:  0.6,
		Tool:        "axe",
		ToolLevel:   0,
		Flammable:   true,
		CanBurn:     true,
	})

	// Copper Block
	registerBlock(BlockCopperBlock, "minecraft:copper_block", &BlockProperties{
		Name:        "Copper Block",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  6.0,
		Tool:        "pickaxe",
		ToolLevel:   1,
	})

	// Amethyst Block
	registerBlock(BlockAmethystBlock, "minecraft:amethyst_block", &BlockProperties{
		Name:        "Amethyst Block",
		Solid:       true,
		Hardness:    1.5,
		Resistance:  1.5,
		Tool:        "pickaxe",
		ToolLevel:   1,
	})

	// Budding Amethyst
	registerBlock(BlockBuddingAmethyst, "minecraft:budding_amethyst", &BlockProperties{
		Name:        "Budding Amethyst",
		Solid:       true,
		Hardness:    1.5,
		Resistance:  1.5,
		Tool:        "pickaxe",
		ToolLevel:   1,
	})

	// Sculk
	registerBlock(BlockSculk, "minecraft:sculk", &BlockProperties{
		Name:        "Sculk",
		Solid:       true,
		Hardness:    1.5,
		Resistance:  1.5,
		Tool:        "hoe",
		ToolLevel:   0,
	})

	// Sculk Sensor
	registerBlock(BlockSculkSensor, "minecraft:sculk_sensor", &BlockProperties{
		Name:        "Sculk Sensor",
		Solid:       true,
		Hardness:    1.5,
		Resistance:  1.5,
		Tool:        "hoe",
		ToolLevel:   0,
	})

	// Sculk Catalyst
	registerBlock(BlockSculkCatalyst, "minecraft:sculk_catalyst", &BlockProperties{
		Name:        "Sculk Catalyst",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "hoe",
		ToolLevel:   0,
	})

	// Sculk Shrieker
	registerBlock(BlockSculkShrieker, "minecraft:sculk_shrieker", &BlockProperties{
		Name:        "Sculk Shrieker",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "hoe",
		ToolLevel:   0,
	})

	// Reinforced Deepslate
	registerBlock(BlockReinforcedDeepslate, "minecraft:reinforced_deepslate", &BlockProperties{
		Name:        "Reinforced Deepslate",
		Solid:       true,
		Hardness:    55.0,
		Resistance:  3600000.0,
		Tool:        "pickaxe",
		ToolLevel:   3,
	})

	// Netherrack
	registerBlock(BlockNetherrack, "minecraft:netherrack", &BlockProperties{
		Name:        "Netherrack",
		Solid:       true,
		Hardness:    0.4,
		Resistance:  0.4,
		Tool:        "pickaxe",
		ToolLevel:   0,
	})

	// End Stone
	registerBlock(BlockEndStone, "minecraft:end_stone", &BlockProperties{
		Name:        "End Stone",
		Solid:       true,
		Hardness:    3.0,
		Resistance:  3.0,
		Tool:        "shovel",
		ToolLevel:   0,
	})

	// Netherite Block
	registerBlock(BlockNetheriteBlock, "minecraft:netherite_block", &BlockProperties{
		Name:        "Netherite Block",
		Solid:       true,
		Hardness:    50.0,
		Resistance:  1200.0,
		Tool:        "pickaxe",
		ToolLevel:   4,
	})

	blockCount = len(blockRegistry)
}

// registerBlock registers a block in the registry.
func registerBlock(id BlockID, name string, props *BlockProperties) {
	props.ID = id
	blockRegistry[id] = props
	blockRegistryByName[name] = id
	blockStateRegistry[id] = &BlockState{
		ID:       id,
		Name:     name,
		Properties: make(map[string]interface{}),
	}
}

// GetBlockProperties returns block properties by ID.
func GetBlockProperties(id BlockID) (*BlockProperties, bool) {
	blockMutex.RLock()
	defer blockMutex.RUnlock()

	props, ok := blockRegistry[id]
	return props, ok
}

// GetBlockByName returns block ID by name.
func GetBlockByName(name string) (BlockID, bool) {
	blockMutex.RLock()
	defer blockMutex.RUnlock()

	id, ok := blockRegistryByName[name]
	return id, ok
}

// GetBlockName returns the name of a block by ID.
func GetBlockName(id BlockID) (string, bool) {
	blockMutex.RLock()
	defer blockMutex.RUnlock()

	if props, ok := blockRegistry[id]; ok {
		return props.Name, true
	}

	return "", false
}

// BlockCount returns the number of registered blocks.
func BlockCount() int {
	blockMutex.RLock()
	defer blockMutex.RUnlock()

	return blockCount
}

// GetAllBlocks returns all registered blocks.
func GetAllBlocks() map[BlockID]*BlockProperties {
	blockMutex.RLock()
	defer blockMutex.RUnlock()

	result := make(map[BlockID]*BlockProperties, len(blockRegistry))
	for id, props := range blockRegistry {
		result[id] = props
	}

	return result
}

// IsSolid checks if a block is solid.
func IsSolid(id BlockID) bool {
	if props, ok := GetBlockProperties(id); ok {
		return props.Solid
	}
	return false
}

// IsTransparent checks if a block is transparent.
func IsTransparent(id BlockID) bool {
	if props, ok := GetBlockProperties(id); ok {
		return props.Transparent
	}
	return false
}

// GetHardness returns the hardness of a block.
func GetHardness(id BlockID) float32 {
	if props, ok := GetBlockProperties(id); ok {
		return props.Hardness
	}
	return 0
}

// GetLightLevel returns the light level of a block.
func GetLightLevel(id BlockID) int {
	if props, ok := GetBlockProperties(id); ok {
		return props.LightLevel
	}
	return 0
}
