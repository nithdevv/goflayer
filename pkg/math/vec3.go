// Package math предоставляет математические утилиты для goflayer.
//
// Этот пакет содержит 3D векторную математику и другие математические функции,
// необходимые для работы с Minecraft миром.
package math

import (
	"fmt"
	"math"
)

// Vec3 представляет трехмерный вектор с координатами X, Y, Z.
//
// Vec3 используется extensively во всем goflayer для представления:
// - Позиций сущностей и блоков
// - Векторов движения и скорости
// - Направлений взгляда
// - Расстояний между объектами
//
// Пример использования:
//
//	pos := Vec3{10.5, 64, -3.2}
//	vel := Vec3{0, 0, 1} // Движение вдоль оси Z
//	distance := pos.DistanceTo(&Vec3{0, 0, 0})
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// NewVec3 создает новый Vec3 из заданных координат.
func NewVec3(x, y, z float64) *Vec3 {
	return &Vec3{X: x, Y: y, Z: z}
}

// String возвращает строковое представление вектора.
func (v *Vec3) String() string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", v.X, v.Y, v.Z)
}

// Clone создает копию вектора.
func (v *Vec3) Clone() *Vec3 {
	return &Vec3{X: v.X, Y: v.Y, Z: v.Z}
}

// Set устанавливает новые координаты вектора.
func (v *Vec3) Set(x, y, z float64) *Vec3 {
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// Update обновляет координаты вектора из другого вектора.
func (v *Vec3) Update(other *Vec3) *Vec3 {
	v.X = other.X
	v.Y = other.Y
	v.Z = other.Z
	return v
}

// Floor возвращает новый вектор с координатами, округленными вниз.
// Полезно для получения позиции блока из позиции сущности.
func (v *Vec3) Floor() *Vec3 {
	return &Vec3{
		X: math.Floor(v.X),
		Y: math.Floor(v.Y),
		Z: math.Floor(v.Z),
	}
}

// Ceil возвращает новый вектор с координатами, округленными вверх.
func (v *Vec3) Ceil() *Vec3 {
	return &Vec3{
		X: math.Ceil(v.X),
		Y: math.Ceil(v.Y),
		Z: math.Ceil(v.Z),
	}
}

// Round возвращает новый вектор с округленными координатами.
func (v *Vec3) Round() *Vec3 {
	return &Vec3{
		X: math.Round(v.X),
		Y: math.Round(v.Y),
		Z: math.Round(v.Z),
	}
}

// Abs возвращает новый вектор с абсолютными значениями координат.
func (v *Vec3) Abs() *Vec3 {
	return &Vec3{
		X: math.Abs(v.X),
		Y: math.Abs(v.Y),
		Z: math.Abs(v.Z),
	}
}

// Add прибавляет другой вектор к этому вектору.
// Возвращает этот же вектор для chaining.
func (v *Vec3) Add(other *Vec3) *Vec3 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// Sub вычитает другой вектор из этого вектора.
// Возвращает этот же вектор для chaining.
func (v *Vec3) Sub(other *Vec3) *Vec3 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// Scaled создает новый вектор, масштабированный на заданный коэффициент.
func (v *Vec3) Scaled(scale float64) *Vec3 {
	return &Vec3{
		X: v.X * scale,
		Y: v.Y * scale,
		Z: v.Z * scale,
	}
}

// Scale масштабирует этот вектор на заданный коэффициент.
// Возвращает этот же вектор для chaining.
func (v *Vec3) Scale(scale float64) *Vec3 {
	v.X *= scale
	v.Y *= scale
	v.Z *= scale
	return v
}

// Plus создает новый вектор, равный сумме этого и другого вектора.
func (v *Vec3) Plus(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Minus создает новый вектор, равный разности этого и другого вектора.
func (v *Vec3) Minus(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Negative возвращает новый вектор с инвертированными координатами.
func (v *Vec3) Negative() *Vec3 {
	return &Vec3{
		X: -v.X,
		Y: -v.Y,
		Z: -v.Z,
	}
}

// Negate инвертирует координаты этого вектора.
// Возвращает этот же вектор для chaining.
func (v *Vec3) Negate() *Vec3 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

// Dot вычисляет скалярное произведение с другим вектором.
func (v *Vec3) Dot(other *Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross вычисляет векторное произведение с другим вектором.
// Возвращает новый вектор.
func (v *Vec3) Cross(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// LengthSquared возвращает квадрат длины вектора.
// Быстрее, чем Length(), если нужна только сравнительная оценка.
func (v *Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Length возвращает длину (мagnitude) вектора.
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Normalize возвращает нормализованный (единичный) вектор.
// Если вектор нулевой, возвращает нулевой вектор.
func (v *Vec3) Normalize() *Vec3 {
	len := v.Length()
	if len == 0 {
		return &Vec3{}
	}
	return v.Scaled(1 / len)
}

// DistanceTo вычисляет Euclidean расстояние до другого вектора.
func (v *Vec3) DistanceTo(other *Vec3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// DistanceToSquared вычисляет квадрат расстояния до другого вектора.
// Быстрее, чем DistanceTo(), если нужна только сравнительная оценка.
func (v *Vec3) DistanceToSquared(other *Vec3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

// Equals проверяет, равен ли этот вектор другому с заданной точностью.
func (v *Vec3) Equals(other *Vec3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// ApproxEquals проверяет, примерно ли равен этот вектор другому с заданной погрешностью.
func (v *Vec3) ApproxEquals(other *Vec3, epsilon float64) bool {
	return math.Abs(v.X-other.X) < epsilon &&
		math.Abs(v.Y-other.Y) < epsilon &&
		math.Abs(v.Z-other.Z) < epsilon
}

// Lerp выполняет линейную интерполяцию между этим и другим вектором.
// t должно быть в диапазоне [0, 1].
// Возвращает новый вектор.
func (v *Vec3) Lerp(other *Vec3, t float64) *Vec3 {
	return &Vec3{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
		Z: v.Z + (other.Z-v.Z)*t,
	}
}

// ManhattanDistance возвращает Manhattan расстояние до другого вектора.
// Это сумма абсолютных разностей координат.
func (v *Vec3) ManhattanDistance(other *Vec3) float64 {
	return math.Abs(v.X-other.X) + math.Abs(v.Y-other.Y) + math.Abs(v.Z-other.Z)
}

// ToArray преобразует вектор в массив из 3 элементов.
func (v *Vec3) ToArray() [3]float64 {
	return [3]float64{v.X, v.Y, v.Z}
}

// ToSlice преобразует вектор в слайс.
func (v *Vec3) ToSlice() []float64 {
	return []float64{v.X, v.Y, v.Z}
}

// AxesID возвращает индекс оси с наибольшим абсолютным значением.
// 0 = X, 1 = Y, 2 = Z.
func (v *Vec3) AxesID() int {
	absX := math.Abs(v.X)
	absY := math.Abs(v.Y)
	absZ := math.Abs(v.Z)

	if absX > absY {
		if absX > absZ {
			return 0 // X
		}
		return 2 // Z
	}
	if absY > absZ {
		return 1 // Y
	}
	return 2 // Z
}

// Offset возвращает новый вектор, смещенный на заданные значения.
// Это удобная обертка вокруг Plus.
func (v *Vec3) Offset(dx, dy, dz float64) *Vec3 {
	return &Vec3{
		X: v.X + dx,
		Y: v.Y + dy,
		Z: v.Z + dz,
	}
}

// Zero проверяет, является ли вектор нулевым.
func (v *Vec3) Zero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

// Yaw вычисляет угол yaw (горизонтальное направление) в радианах.
// Yaw - это угол в горизонтальной плоскости XZ.
func (v *Vec3) Yaw() float64 {
	return math.Atan2(v.Z, v.X)
}

// Pitch вычисляет угол pitch (вертикальное направление) в радианах.
// Pitch - это угол относительно горизонтали.
func (v *Vec3) Pitch() float64 {
	len := math.Sqrt(v.X*v.X + v.Z*v.Z)
	return math.Atan2(v.Y, len)
}

// YawPitch возвращает yaw и pitch углы в радианах.
func (v *Vec3) YawPitch() (float64, float64) {
	yaw := math.Atan2(v.Z, v.X)
	len := math.Sqrt(v.X*v.X + v.Z*v.Z)
	pitch := math.Atan2(v.Y, len)
	return yaw, pitch
}

// FromYawPitch создает новый вектор направления из yaw и pitch углов в радианах.
// Это обратная операция к YawPitch().
func FromYawPitch(yaw, pitch float64) *Vec3 {
	cosPitch := math.Cos(pitch)
	sinPitch := math.Sin(pitch)
	cosYaw := math.Cos(yaw)
	sinYaw := math.Sin(yaw)

	return &Vec3{
		X: -cosYaw * cosPitch,
		Y: -sinPitch,
		Z: -sinYaw * cosPitch,
	}
}

// ViewDirection создает вектор направления взгляда из yaw и pitch углов.
// Это альтернативное название для FromYawPitch для совместимости с mineflayer API.
func ViewDirection(yaw, pitch float64) *Vec3 {
	return FromYawPitch(yaw, pitch)
}

// Clamp ограничивает координаты вектора в заданном диапазоне.
// Возвращает новый вектор.
func (v *Vec3) Clamp(min, max float64) *Vec3 {
	return &Vec3{
		X: clamp(min, v.X, max),
		Y: clamp(min, v.Y, max),
		Z: clamp(min, v.Z, max),
	}
}

// ClampEach ограничивает каждую координату в своем диапазоне.
// Возвращает новый вектор.
func (v *Vec3) ClampEach(minX, maxX, minY, maxY, minZ, maxZ float64) *Vec3 {
	return &Vec3{
		X: clamp(minX, v.X, maxX),
		Y: clamp(minY, v.Y, maxY),
		Z: clamp(minZ, v.Z, maxZ),
	}
}

// Reflect отражает вектор относительно нормали.
// Возвращает новый вектор.
func (v *Vec3) Reflect(normal *Vec3) *Vec3 {
	// r = v - 2 * (v . n) * n
	dot := v.Dot(normal)
	return v.Minus(normal.Scaled(2 * dot))
}

// Project проектирует этот вектор на другой вектор.
// Возвращает новый вектор.
func (v *Vec3) Project(other *Vec3) *Vec3 {
	// proj = (v . other / |other|^2) * other
	lenSquared := other.LengthSquared()
	if lenSquared == 0 {
		return &Vec3{}
	}
	dot := v.Dot(other)
	return other.Scaled(dot / lenSquared)
}

// Clamp ограничивает значение x в диапазоне [min, max].
func clamp(min, x, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// EuclideanMod выполняет евклидово деление по модулю.
// В отличие от стандартного % в Go, результат всегда положительный.
func EuclideanMod(numerator, denominator float64) float64 {
	result := math.Mod(numerator, denominator)
	if result < 0 {
		result += denominator
	}
	return result
}

// DegreesToRadians преобразует градусы в радианы.
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RadiansToDegrees преобразует радианы в градусы.
func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}
