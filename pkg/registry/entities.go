// Package registry provides Minecraft 1.20.1 entity type registry.
package registry

import (
	"sync"
)

// EntityType represents an entity type.
type EntityType int32

// EntityProperties holds entity properties.
type EntityProperties struct {
	ID          EntityType
	Name        string
	Category    EntityCategory
	Width       float32
	Height      float32
	Health      float32
	Hostile     bool
	Passive     bool
	Tamable     bool
	Breedable   bool
	FireImmune  bool
	WaterImmune bool
}

// EntityCategory represents the category of an entity.
type EntityCategory string

const (
	EntityCategoryMonster    EntityCategory = "monster"
	EntityCategoryCreature   EntityCategory = "creature"
	EntityCategoryAmbient    EntityCategory = "ambient"
	EntityCategoryWater      EntityCategory = "water"
	EntityCategoryMisc       EntityCategory = "misc"
	EntityCategoryNPC        EntityCategory = "npc"
	EntityCategoryProjectile EntityCategory = "projectile"
	EntityCategoryItem       EntityCategory = "item"
	EntityCategoryPlayer     EntityCategory = "player"
)

var (
	entityRegistry      map[EntityType]*EntityProperties
	entityRegistryByName map[string]EntityType
	entityMutex         sync.RWMutex
	entityCount         int
)

func init() {
	initializeEntities()
}

// Entity type constants for Minecraft 1.20.1
const (
	EntityTypeAllay EntityType = 0
	EntityTypeAreaEffectCloud EntityType = 1
	EntityTypeArmorStand EntityType = 2
	EntityTypeArrow EntityType = 3
	EntityTypeAxolotl EntityType = 4
	EntityTypeBat EntityType = 5
	EntityTypeBee EntityType = 6
	EntityTypeBlaze EntityType = 7
	EntityTypeBlockDisplay EntityType = 8
	EntityTypeBoat EntityType = 9
	EntityTypeCamel EntityType = 10
	EntityTypeCat EntityType = 11
	EntityTypeCaveSpider EntityType = 12
	EntityTypeChestBoat EntityType = 13
	EntityTypeChestMinecart EntityType = 14
	EntityTypeChicken EntityType = 15
	EntityTypeCod EntityType = 16
	EntityTypeCommandBlockMinecart EntityType = 17
	EntityTypeCreeper EntityType = 18
	EntityTypeDolphin EntityType = 19
	EntityTypeDonkey EntityType = 20
	EntityTypeDragonFireball EntityType = 21
	EntityTypeDrowned EntityType = 22
	EntityTypeEgg EntityType = 23
	EntityTypeElderGuardian EntityType = 24
	EntityTypeEndCrystal EntityType = 25
	EntityTypeEnderDragon EntityType = 26
	EntityTypeEnderman EntityType = 27
	EntityTypeEndermite EntityType = 28
	EntityTypeEvoker EntityType = 29
	EntityTypeEvokerFangs EntityType = 30
	EntityTypeExperienceBottle EntityType = 31
	EntityTypeExperienceOrb EntityType = 32
	EntityTypeEyeOfEnder EntityType = 33
	EntityTypeFallingBlock EntityType = 34
	EntityTypeFireworkRocket EntityType = 35
	EntityTypeFireworkRocketEntity EntityType = 36
	EntityTypeFox EntityType = 37
	EntityTypeFrog EntityType = 38
	EntityTypeGhast EntityType = 39
	EntityTypeGiant EntityType = 40
	EntityTypeGlowItemFrame EntityType = 41
	EntityTypeGlowSquid EntityType = 42
	EntityTypeGoat EntityType = 43
	EntityTypeGuardian EntityType = 44
	EntityTypeHoglin EntityType = 45
	EntityTypeHopperMinecart EntityType = 46
	EntityTypeHorse EntityType = 47
	EntityTypeHusk EntityType = 48
	EntityTypeIllusioner EntityType = 49
	EntityTypeInteraction EntityType = 50
	EntityTypeIronGolem EntityType = 51
	EntityTypeItem EntityType = 52
	EntityTypeItemDisplay EntityType = 53
	EntityTypeItemFrame EntityType = 54
	EntityTypeFireball EntityType = 55
	EntityTypeLeashKnot EntityType = 56
	EntityTypeLightningBolt EntityType = 57
	EntityTypeLlama EntityType = 58
	EntityTypeLlamaSpit EntityType = 59
	EntityTypeMagmaCube EntityType = 60
	EntityTypeMarker EntityType = 61
	EntityTypeMinecart EntityType = 62
	EntityTypeMooshroom EntityType = 63
	EntityTypeOcelot EntityType = 64
	EntityTypePainting EntityType = 65
	EntityTypePanda EntityType = 66
	EntityTypeParrot EntityType = 67
	EntityTypePhantom EntityType = 68
	EntityTypePig EntityType = 69
	EntityTypePiglin EntityType = 70
	EntityTypePiglinBrute EntityType = 71
	EntityTypePillager EntityType = 72
	EntityTypePolarBear EntityType = 73
	EntityTypePufferfish EntityType = 74
	EntityTypeRabbit EntityType = 75
	EntityTypeRavager EntityType = 76
	EntityTypeSalmon EntityType = 77
	EntityTypeSheep EntityType = 78
	EntityTypeShulker EntityType = 79
	EntityTypeShulkerBullet EntityType = 80
	EntityTypeSilverfish EntityType = 81
	EntityTypeSkeleton EntityType = 82
	EntityTypeSkeletonHorse EntityType = 83
	EntityTypeSlime EntityType = 84
	EntityTypeSmallFireball EntityType = 85
	EntityTypeSniffer EntityType = 86
	EntityTypeSnowGolem EntityType = 87
	EntityTypeSnowball EntityType = 88
	EntityTypeSpawnerMinecart EntityType = 89
	EntityTypeSpectralArrow EntityType = 90
	EntityTypeSpider EntityType = 91
	EntityTypeSquid EntityType = 92
	EntityTypeStray EntityType = 93
	EntityTypeStrider EntityType = 94
	EntityTypeTadpole EntityType = 95
	EntityTypeTextDisplay EntityType = 96
	EntityTypeThrownEgg EntityType = 97
	EntityTypeThrownEnderpearl EntityType = 98
	EntityTypeThrownPotion EntityType = 99
	EntityTypeThrownExpBottle EntityType = 100
	EntityTypeTraderLlama EntityType = 101
	EntityTypeTrident EntityType = 102
	EntityTypeTropicalFish EntityType = 103
	EntityTypeTurtle EntityType = 104
	EntityTypeVex EntityType = 105
	EntityTypeVillager EntityType = 106
	EntityTypeVindicator EntityType = 107
	EntityTypeWanderingTrader EntityType = 108
	EntityTypeWarden EntityType = 109
	EntityTypeWindCharge EntityType = 110
	EntityTypeWitch EntityType = 111
	EntityTypeWither EntityType = 112
	EntityTypeWitherSkeleton EntityType = 113
	EntityTypeWitherSkull EntityType = 114
	EntityTypeWolf EntityType = 115
	EntityTypeZoglin EntityType = 116
	EntityTypeZombie EntityType = 117
	EntityTypeZombieHorse EntityType = 118
	EntityTypeZombieVillager EntityType = 119
	EntityTypeZombifiedPiglin EntityType = 120
	EntityTypePlayer EntityType = 121
	EntityTypeFishingBobber EntityType = 122
)

// initializeEntities initializes the entity registry.
func initializeEntities() {
	entityRegistry = make(map[EntityType]*EntityProperties)
	entityRegistryByName = make(map[string]EntityType)

	// Allay
	registerEntity(EntityTypeAllay, "minecraft:allay", &EntityProperties{
		Name:       "Allay",
		Category:   EntityCategoryCreature,
		Width:      0.6,
		Height:     0.6,
		Health:     20,
		Hostile:    false,
		Passive:    true,
		FireImmune: false,
	})

	// Area Effect Cloud
	registerEntity(EntityTypeAreaEffectCloud, "minecraft:area_effect_cloud", &EntityProperties{
		Name:       "Area Effect Cloud",
		Category:   EntityCategoryMisc,
		Width:      6.0,
		Height:     0.5,
		Health:     0, // Invulnerable
		Hostile:    false,
		Passive:    false,
	})

	// Armor Stand
	registerEntity(EntityTypeArmorStand, "minecraft:armor_stand", &EntityProperties{
		Name:       "Armor Stand",
		Category:   EntityCategoryMisc,
		Width:      0.5,
		Height:     1.975,
		Health:     25,
		Hostile:    false,
		Passive:    true,
	})

	// Arrow
	registerEntity(EntityTypeArrow, "minecraft:arrow", &EntityProperties{
		Name:       "Arrow",
		Category:   EntityCategoryProjectile,
		Width:      0.5,
		Height:     0.5,
		Health:     0,
		Hostile:    false,
	})

	// Axolotl
	registerEntity(EntityTypeAxolotl, "minecraft:axolotl", &EntityProperties{
		Name:         "Axolotl",
		Category:     EntityCategoryWater,
		Width:        1.3,
		Height:       0.6,
		Health:       14,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Bat
	registerEntity(EntityTypeBat, "minecraft:bat", &EntityProperties{
		Name:       "Bat",
		Category:   EntityCategoryAmbient,
		Width:      0.5,
		Height:     0.9,
		Health:     6,
		Hostile:    false,
		Passive:    true,
	})

	// Bee
	registerEntity(EntityTypeBee, "minecraft:bee", &EntityProperties{
		Name:       "Bee",
		Category:   EntityCategoryCreature,
		Width:      0.7,
		Height:     0.6,
		Health:     10,
		Hostile:    false,
		Passive:    true,
	})

	// Blaze
	registerEntity(EntityTypeBlaze, "minecraft:blaze", &EntityProperties{
		Name:       "Blaze",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.8,
		Health:     20,
		Hostile:    true,
		Passive:    false,
		FireImmune: true,
	})

	// Boat
	registerEntity(EntityTypeBoat, "minecraft:boat", &EntityProperties{
		Name:       "Boat",
		Category:   EntityCategoryMisc,
		Width:      1.375,
		Height:     0.5625,
		Health:     0,
		Hostile:    false,
	})

	// Camel
	registerEntity(EntityTypeCamel, "minecraft:camel", &EntityProperties{
		Name:       "Camel",
		Category:   EntityCategoryCreature,
		Width:      1.7,
		Height:     2.375,
		Health:     32,
		Hostile:    false,
		Passive:    true,
		Tamable:    true,
	})

	// Cat
	registerEntity(EntityTypeCat, "minecraft:cat", &EntityProperties{
		Name:       "Cat",
		Category:   EntityCategoryCreature,
		Width:      0.6,
		Height:     0.7,
		Health:     10,
		Hostile:    false,
		Passive:    true,
		Tamable:    true,
	})

	// Cave Spider
	registerEntity(EntityTypeCaveSpider, "minecraft:cave_spider", &EntityProperties{
		Name:       "Cave Spider",
		Category:   EntityCategoryMonster,
		Width:      0.7,
		Height:     0.5,
		Health:     12,
		Hostile:    true,
		Passive:    false,
	})

	// Chicken
	registerEntity(EntityTypeChicken, "minecraft:chicken", &EntityProperties{
		Name:       "Chicken",
		Category:   EntityCategoryCreature,
		Width:      0.4,
		Height:     0.7,
		Health:     4,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Cod
	registerEntity(EntityTypeCod, "minecraft:cod", &EntityProperties{
		Name:         "Cod",
		Category:     EntityCategoryWater,
		Width:        0.5,
		Height:       0.3,
		Health:       3,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Creeper
	registerEntity(EntityTypeCreeper, "minecraft:creeper", &EntityProperties{
		Name:       "Creeper",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.7,
		Health:     20,
		Hostile:    true,
		Passive:    false,
	})

	// Dolphin
	registerEntity(EntityTypeDolphin, "minecraft:dolphin", &EntityProperties{
		Name:         "Dolphin",
		Category:     EntityCategoryWater,
		Width:        0.9,
		Height:       0.6,
		Health:       10,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Donkey
	registerEntity(EntityTypeDonkey, "minecraft:donkey", &EntityProperties{
		Name:       "Donkey",
		Category:   EntityCategoryCreature,
		Width:      1.3964844,
		Height:     1.5,
		Health:     15, // 23-29
		Hostile:    false,
		Passive:    true,
		Tamable:    true,
	})

	// Drowned
	registerEntity(EntityTypeDrowned, "minecraft:drowned", &EntityProperties{
		Name:         "Drowned",
		Category:     EntityCategoryMonster,
		Width:        0.6,
		Height:       1.95,
		Health:       20,
		Hostile:      true,
		Passive:      false,
		WaterImmune:  true,
	})

	// Elder Guardian
	registerEntity(EntityTypeElderGuardian, "minecraft:elder_guardian", &EntityProperties{
		Name:         "Elder Guardian",
		Category:     EntityCategoryMonster,
		Width:        1.9975,
		Height:       1.9975,
		Health:       80,
		Hostile:      true,
		Passive:      false,
		WaterImmune:  true,
	})

	// Ender Dragon
	registerEntity(EntityTypeEnderDragon, "minecraft:ender_dragon", &EntityProperties{
		Name:       "Ender Dragon",
		Category:   EntityCategoryMonster,
		Width:      16.0,
		Height:     8.0,
		Health:     200,
		Hostile:    true,
		Passive:    false,
	})

	// Enderman
	registerEntity(EntityTypeEnderman, "minecraft:enderman", &EntityProperties{
		Name:       "Enderman",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     2.9,
		Health:     40,
		Hostile:    true,
		Passive:    false,
	})

	// Evoker
	registerEntity(EntityTypeEvoker, "minecraft:evoker", &EntityProperties{
		Name:       "Evoker",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     24,
		Hostile:    true,
		Passive:    false,
	})

	// Experience Orb
	registerEntity(EntityTypeExperienceOrb, "minecraft:experience_orb", &EntityProperties{
		Name:     "Experience Orb",
		Category: EntityCategoryMisc,
		Width:    0.5,
		Height:   0.5,
		Health:   5,
	})

	// Fox
	registerEntity(EntityTypeFox, "minecraft:fox", &EntityProperties{
		Name:       "Fox",
		Category:   EntityCategoryCreature,
		Width:      0.6,
		Height:     0.7,
		Health:     10,
		Hostile:    false,
		Passive:    true,
	})

	// Frog
	registerEntity(EntityTypeFrog, "minecraft:frog", &EntityProperties{
		Name:       "Frog",
		Category:   EntityCategoryCreature,
		Width:      0.9,
		Height:     0.5,
		Health:     10,
		Hostile:    false,
		Passive:    true,
	})

	// Ghast
	registerEntity(EntityTypeGhast, "minecraft:ghast", &EntityProperties{
		Name:       "Ghast",
		Category:   EntityCategoryMonster,
		Width:      4.0,
		Height:     4.0,
		Health:     10,
		Hostile:    true,
		Passive:    false,
		FireImmune: true,
	})

	// Giant
	registerEntity(EntityTypeGiant, "minecraft:giant", &EntityProperties{
		Name:       "Giant",
		Category:   EntityCategoryMonster,
		Width:      3.6,
		Height:     12.0,
		Health:     100,
		Hostile:    true,
		Passive:    false,
	})

	// Glow Squid
	registerEntity(EntityTypeGlowSquid, "minecraft:glow_squid", &EntityProperties{
		Name:         "Glow Squid",
		Category:     EntityCategoryWater,
		Width:        0.8,
		Height:       0.8,
		Health:       10,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Goat
	registerEntity(EntityTypeGoat, "minecraft:goat", &EntityProperties{
		Name:       "Goat",
		Category:   EntityCategoryCreature,
		Width:      1.3,
		Height:     1.45,
		Health:     10, // 5-7
		Hostile:    false,
		Passive:    true,
	})

	// Guardian
	registerEntity(EntityTypeGuardian, "minecraft:guardian", &EntityProperties{
		Name:         "Guardian",
		Category:     EntityCategoryMonster,
		Width:        0.85,
		Height:       0.85,
		Health:       30,
		Hostile:      true,
		Passive:      false,
		WaterImmune:  true,
	})

	// Hoglin
	registerEntity(EntityTypeHoglin, "minecraft:hoglin", &EntityProperties{
		Name:       "Hoglin",
		Category:   EntityCategoryCreature,
		Width:      1.3964844,
		Height:     1.4,
		Health:     40,
		Hostile:    true,
		Passive:    false,
		Breedable:  true,
		FireImmune: true,
	})

	// Horse
	registerEntity(EntityTypeHorse, "minecraft:horse", &EntityProperties{
		Name:       "Horse",
		Category:   EntityCategoryCreature,
		Width:      1.3964844,
		Height:     1.6,
		Health:     15, // 20-30
		Hostile:    false,
		Passive:    true,
		Tamable:    true,
	})

	// Husk
	registerEntity(EntityTypeHusk, "minecraft:husk", &EntityProperties{
		Name:       "Husk",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     20,
		Hostile:    true,
		Passive:    false,
	})

	// Iron Golem
	registerEntity(EntityTypeIronGolem, "minecraft:iron_golem", &EntityProperties{
		Name:       "Iron Golem",
		Category:   EntityCategoryCreature,
		Width:      1.4,
		Height:     2.7,
		Health:     100,
		Hostile:    false,
		Passive:    false,
	})

	// Item
	registerEntity(EntityTypeItem, "minecraft:item", &EntityProperties{
		Name:     "Item",
		Category: EntityCategoryItem,
		Width:    0.25,
		Height:   0.25,
		Health:   5,
	})

	// Fireball
	registerEntity(EntityTypeFireball, "minecraft:fireball", &EntityProperties{
		Name:       "Fireball",
		Category:   EntityCategoryProjectile,
		Width:      0.3125,
		Height:     0.3125,
		Health:     1,
		Hostile:    false,
		FireImmune: true,
	})

	// Llama
	registerEntity(EntityTypeLlama, "minecraft:llama", &EntityProperties{
		Name:       "Llama",
		Category:   EntityCategoryCreature,
		Width:      0.9,
		Height:     1.87,
		Health:     15, // 22-33
		Hostile:    false,
		Passive:    true,
		Tamable:    true,
		Breedable:  true,
	})

	// Magma Cube
	registerEntity(EntityTypeMagmaCube, "minecraft:magma_cube", &EntityProperties{
		Name:       "Magma Cube",
		Category:   EntityCategoryMonster,
		Width:      2.04,
		Height:     2.04,
		Health:     16,
		Hostile:    true,
		Passive:    false,
		FireImmune: true,
	})

	// Mooshroom
	registerEntity(EntityTypeMooshroom, "minecraft:mooshroom", &EntityProperties{
		Name:       "Mooshroom",
		Category:   EntityCategoryCreature,
		Width:      0.9,
		Height:     1.4,
		Health:     10,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Ocelot
	registerEntity(EntityTypeOcelot, "minecraft:ocelot", &EntityProperties{
		Name:       "Ocelot",
		Category:   EntityCategoryCreature,
		Width:      0.6,
		Height:     0.7,
		Health:     10,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Painting
	registerEntity(EntityTypePainting, "minecraft:painting", &EntityProperties{
		Name:     "Painting",
		Category: EntityCategoryMisc,
		Width:    0.5,
		Height:   0.5,
		Health:   0,
	})

	// Panda
	registerEntity(EntityTypePanda, "minecraft:panda", &EntityProperties{
		Name:       "Panda",
		Category:   EntityCategoryCreature,
		Width:      1.2,
		Height:     1.25,
		Health:     20,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Parrot
	registerEntity(EntityTypeParrot, "minecraft:parrot", &EntityProperties{
		Name:       "Parrot",
		Category:   EntityCategoryCreature,
		Width:      0.5,
		Height:     0.9,
		Health:     6,
		Hostile:    false,
		Passive:    true,
		Tamable:    true,
	})

	// Phantom
	registerEntity(EntityTypePhantom, "minecraft:phantom", &EntityProperties{
		Name:       "Phantom",
		Category:   EntityCategoryMonster,
		Width:      0.9,
		Height:     0.5,
		Health:     20,
		Hostile:    true,
		Passive:    false,
	})

	// Pig
	registerEntity(EntityTypePig, "minecraft:pig", &EntityProperties{
		Name:       "Pig",
		Category:   EntityCategoryCreature,
		Width:      0.9,
		Height:     0.9,
		Health:     10,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Piglin
	registerEntity(EntityTypePiglin, "minecraft:piglin", &EntityProperties{
		Name:       "Piglin",
		Category:   EntityCategoryCreature,
		Width:      0.6,
		Height:     1.95,
		Health:     16,
		Hostile:    true,
		Passive:    false,
		FireImmune: true,
	})

	// Piglin Brute
	registerEntity(EntityTypePiglinBrute, "minecraft:piglin_brute", &EntityProperties{
		Name:       "Piglin Brute",
		Category:   EntityCategoryCreature,
		Width:      0.6,
		Height:     1.95,
		Health:     50,
		Hostile:    true,
		Passive:    false,
		FireImmune: true,
	})

	// Pillager
	registerEntity(EntityTypePillager, "minecraft:pillager", &EntityProperties{
		Name:       "Pillager",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     24,
		Hostile:    true,
		Passive:    false,
	})

	// Polar Bear
	registerEntity(EntityTypePolarBear, "minecraft:polar_bear", &EntityProperties{
		Name:       "Polar Bear",
		Category:   EntityCategoryCreature,
		Width:      1.4,
		Height:     1.4,
		Health:     30,
		Hostile:    false,
		Passive:    true,
	})

	// Pufferfish
	registerEntity(EntityTypePufferfish, "minecraft:pufferfish", &EntityProperties{
		Name:         "Pufferfish",
		Category:     EntityCategoryWater,
		Width:        0.7,
		Height:       0.7,
		Health:       3,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Rabbit
	registerEntity(EntityTypeRabbit, "minecraft:rabbit", &EntityProperties{
		Name:       "Rabbit",
		Category:   EntityCategoryCreature,
		Width:      0.4,
		Height:     0.5,
		Health:     3,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Ravager
	registerEntity(EntityTypeRavager, "minecraft:ravager", &EntityProperties{
		Name:       "Ravager",
		Category:   EntityCategoryMonster,
		Width:      1.95,
		Height:     2.2,
		Health:     100,
		Hostile:    true,
		Passive:    false,
	})

	// Salmon
	registerEntity(EntityTypeSalmon, "minecraft:salmon", &EntityProperties{
		Name:         "Salmon",
		Category:     EntityCategoryWater,
		Width:        0.7,
		Height:       0.4,
		Health:       3,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Sheep
	registerEntity(EntityTypeSheep, "minecraft:sheep", &EntityProperties{
		Name:       "Sheep",
		Category:   EntityCategoryCreature,
		Width:      0.9,
		Height:     1.3,
		Health:     8,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Shulker
	registerEntity(EntityTypeShulker, "minecraft:shulker", &EntityProperties{
		Name:       "Shulker",
		Category:   EntityCategoryMonster,
		Width:      1.0,
		Height:     1.0,
		Health:     30,
		Hostile:    true,
		Passive:    false,
	})

	// Silverfish
	registerEntity(EntityTypeSilverfish, "minecraft:silverfish", &EntityProperties{
		Name:       "Silverfish",
		Category:   EntityCategoryMonster,
		Width:      0.4,
		Height:     0.3,
		Health:     8,
		Hostile:    true,
		Passive:    false,
	})

	// Skeleton
	registerEntity(EntityTypeSkeleton, "minecraft:skeleton", &EntityProperties{
		Name:       "Skeleton",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.99,
		Health:     20,
		Hostile:    true,
		Passive:    false,
	})

	// Skeleton Horse
	registerEntity(EntityTypeSkeletonHorse, "minecraft:skeleton_horse", &EntityProperties{
		Name:       "Skeleton Horse",
		Category:   EntityCategoryCreature,
		Width:      1.3964844,
		Height:     1.6,
		Health:     15,
		Hostile:    false,
		Passive:    true,
	})

	// Slime
	registerEntity(EntityTypeSlime, "minecraft:slime", &EntityProperties{
		Name:       "Slime",
		Category:   EntityCategoryMonster,
		Width:      2.04,
		Height:     2.04,
		Health:     16,
		Hostile:    true,
		Passive:    false,
	})

	// Sniffer
	registerEntity(EntityTypeSniffer, "minecraft:sniffer", &EntityProperties{
		Name:       "Sniffer",
		Category:   EntityCategoryCreature,
		Width:      1.9,
		Height:     1.75,
		Health:     14,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Snow Golem
	registerEntity(EntityTypeSnowGolem, "minecraft:snow_golem", &EntityProperties{
		Name:       "Snow Golem",
		Category:   EntityCategoryCreature,
		Width:      0.7,
		Height:     1.9,
		Health:     4,
		Hostile:    false,
		Passive:    false,
	})

	// Spider
	registerEntity(EntityTypeSpider, "minecraft:spider", &EntityProperties{
		Name:       "Spider",
		Category:   EntityCategoryMonster,
		Width:      1.4,
		Height:     0.9,
		Health:     16,
		Hostile:    true,
		Passive:    false,
	})

	// Squid
	registerEntity(EntityTypeSquid, "minecraft:squid", &EntityProperties{
		Name:         "Squid",
		Category:     EntityCategoryWater,
		Width:        0.8,
		Height:       0.8,
		Health:       10,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Stray
	registerEntity(EntityTypeStray, "minecraft:stray", &EntityProperties{
		Name:       "Stray",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.99,
		Health:     20,
		Hostile:    true,
		Passive:    false,
	})

	// Strider
	registerEntity(EntityTypeStrider, "minecraft:strider", &EntityProperties{
		Name:       "Strider",
		Category:   EntityCategoryCreature,
		Width:      0.9,
		Height:     1.7,
		Health:     20,
		Hostile:    false,
		Passive:    true,
		FireImmune: true,
	})

	// Tadpole
	registerEntity(EntityTypeTadpole, "minecraft:tadpole", &EntityProperties{
		Name:         "Tadpole",
		Category:     EntityCategoryWater,
		Width:        0.4,
		Height:       0.3,
		Health:       6,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Trader Llama
	registerEntity(EntityTypeTraderLlama, "minecraft:trader_llama", &EntityProperties{
		Name:       "Trader Llama",
		Category:   EntityCategoryCreature,
		Width:      0.9,
		Height:     1.87,
		Health:     15, // 22-33
		Hostile:    false,
		Passive:    true,
	})

	// Tropical Fish
	registerEntity(EntityTypeTropicalFish, "minecraft:tropical_fish", &EntityProperties{
		Name:         "Tropical Fish",
		Category:     EntityCategoryWater,
		Width:        0.5,
		Height:       0.4,
		Health:       3,
		Hostile:      false,
		Passive:      true,
		WaterImmune:  true,
	})

	// Turtle
	registerEntity(EntityTypeTurtle, "minecraft:turtle", &EntityProperties{
		Name:       "Turtle",
		Category:   EntityCategoryCreature,
		Width:      1.2,
		Height:     0.4,
		Health:     30,
		Hostile:    false,
		Passive:    true,
		Breedable:  true,
	})

	// Vex
	registerEntity(EntityTypeVex, "minecraft:vex", &EntityProperties{
		Name:       "Vex",
		Category:   EntityCategoryMonster,
		Width:      0.4,
		Height:     0.8,
		Health:     14,
		Hostile:    true,
		Passive:    false,
	})

	// Villager
	registerEntity(EntityTypeVillager, "minecraft:villager", &EntityProperties{
		Name:       "Villager",
		Category:   EntityCategoryNPC,
		Width:      0.6,
		Height:     1.95,
		Health:     20,
		Hostile:    false,
		Passive:    true,
	})

	// Vindicator
	registerEntity(EntityTypeVindicator, "minecraft:vindicator", &EntityProperties{
		Name:       "Vindicator",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     24,
		Hostile:    true,
		Passive:    false,
	})

	// Wandering Trader
	registerEntity(EntityTypeWanderingTrader, "minecraft:wandering_trader", &EntityProperties{
		Name:       "Wandering Trader",
		Category:   EntityCategoryNPC,
		Width:      0.6,
		Height:     1.95,
		Health:     20,
		Hostile:    false,
		Passive:    true,
	})

	// Warden
	registerEntity(EntityTypeWarden, "minecraft:warden", &EntityProperties{
		Name:       "Warden",
		Category:   EntityCategoryMonster,
		Width:      0.9,
		Height:     2.9,
		Health:     500,
		Hostile:    true,
		Passive:    false,
	})

	// Witch
	registerEntity(EntityTypeWitch, "minecraft:witch", &EntityProperties{
		Name:       "Witch",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     26,
		Hostile:    true,
		Passive:    false,
	})

	// Wither
	registerEntity(EntityTypeWither, "minecraft:wither", &EntityProperties{
		Name:       "Wither",
		Category:   EntityCategoryMonster,
		Width:      1.5,
		Height:     3.5,
		Health:     300,
		Hostile:    true,
		Passive:    false,
	})

	// Wither Skeleton
	registerEntity(EntityTypeWitherSkeleton, "minecraft:wither_skeleton", &EntityProperties{
		Name:       "Wither Skeleton",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     2.4,
		Health:     20,
		Hostile:    true,
		Passive:    false,
		FireImmune: true,
	})

	// Wolf
	registerEntity(EntityTypeWolf, "minecraft:wolf", &EntityProperties{
		Name:       "Wolf",
		Category:   EntityCategoryCreature,
		Width:      0.6,
		Height:     0.85,
		Health:     8,
		Hostile:    false,
		Passive:    true,
		Tamable:    true,
		Breedable:  true,
	})

	// Zoglin
	registerEntity(EntityTypeZoglin, "minecraft:zoglin", &EntityProperties{
		Name:       "Zoglin",
		Category:   EntityCategoryMonster,
		Width:      1.3964844,
		Height:     1.4,
		Health:     40,
		Hostile:    true,
		Passive:    false,
		FireImmune: true,
	})

	// Zombie
	registerEntity(EntityTypeZombie, "minecraft:zombie", &EntityProperties{
		Name:       "Zombie",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     20,
		Hostile:    true,
		Passive:    false,
	})

	// Zombie Horse
	registerEntity(EntityTypeZombieHorse, "minecraft:zombie_horse", &EntityProperties{
		Name:       "Zombie Horse",
		Category:   EntityCategoryCreature,
		Width:      1.3964844,
		Height:     1.6,
		Health:     15,
		Hostile:    false,
		Passive:    true,
	})

	// Zombie Villager
	registerEntity(EntityTypeZombieVillager, "minecraft:zombie_villager", &EntityProperties{
		Name:       "Zombie Villager",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     20,
		Hostile:    true,
		Passive:    false,
	})

	// Zombified Piglin
	registerEntity(EntityTypeZombifiedPiglin, "minecraft:zombified_piglin", &EntityProperties{
		Name:       "Zombified Piglin",
		Category:   EntityCategoryMonster,
		Width:      0.6,
		Height:     1.95,
		Health:     20,
		Hostile:    false,
		Passive:    false,
		FireImmune: true,
	})

	// Player
	registerEntity(EntityTypePlayer, "minecraft:player", &EntityProperties{
		Name:       "Player",
		Category:   EntityCategoryPlayer,
		Width:      0.6,
		Height:     1.8,
		Health:     20,
		Hostile:    false,
		Passive:    false,
	})

	// Fishing Bobber
	registerEntity(EntityTypeFishingBobber, "minecraft:fishing_bobber", &EntityProperties{
		Name:     "Fishing Bobber",
		Category: EntityCategoryMisc,
		Width:    0.25,
		Height:   0.25,
		Health:   0,
	})

	entityCount = len(entityRegistry)
}

// registerEntity registers an entity in the registry.
func registerEntity(id EntityType, name string, props *EntityProperties) {
	props.ID = id
	entityRegistry[id] = props
	entityRegistryByName[name] = id
}

// GetEntityProperties returns entity properties by type.
func GetEntityProperties(entityType EntityType) (*EntityProperties, bool) {
	entityMutex.RLock()
	defer entityMutex.RUnlock()

	props, ok := entityRegistry[entityType]
	return props, ok
}

// GetEntityByName returns entity type by name.
func GetEntityByName(name string) (EntityType, bool) {
	entityMutex.RLock()
	defer entityMutex.RUnlock()

	entityType, ok := entityRegistryByName[name]
	return entityType, ok
}

// GetEntityName returns the name of an entity by type.
func GetEntityName(entityType EntityType) (string, bool) {
	entityMutex.RLock()
	defer entityMutex.RUnlock()

	if props, ok := entityRegistry[entityType]; ok {
		return props.Name, true
	}

	return "", false
}

// EntityCount returns the number of registered entities.
func EntityCount() int {
	entityMutex.RLock()
	defer entityMutex.RUnlock()

	return entityCount
}

// GetAllEntities returns all registered entities.
func GetAllEntities() map[EntityType]*EntityProperties {
	entityMutex.RLock()
	defer entityMutex.RUnlock()

	result := make(map[EntityType]*EntityProperties, len(entityRegistry))
	for id, props := range entityRegistry {
		result[id] = props
	}

	return result
}

// IsHostile checks if an entity type is hostile.
func IsHostile(entityType EntityType) bool {
	if props, ok := GetEntityProperties(entityType); ok {
		return props.Hostile
	}
	return false
}

// IsPassive checks if an entity type is passive.
func IsPassive(entityType EntityType) bool {
	if props, ok := GetEntityProperties(entityType); ok {
		return props.Passive
	}
	return false
}

// GetEntityHealth returns the base health for an entity type.
func GetEntityHealth(entityType EntityType) float32 {
	if props, ok := GetEntityProperties(entityType); ok {
		return props.Health
	}
	return 0
}

// GetEntitySize returns the width and height for an entity type.
func GetEntitySize(entityType EntityType) (width, height float32) {
	if props, ok := GetEntityProperties(entityType); ok {
		return props.Width, props.Height
	}
	return 0, 0
}

// GetEntityCategory returns the category for an entity type.
func GetEntityCategory(entityType EntityType) EntityCategory {
	if props, ok := GetEntityProperties(entityType); ok {
		return props.Category
	}
	return EntityCategoryMisc
}
