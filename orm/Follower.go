package orm

type Follow struct {
	ID         uint `gorm:"primaryKey;autoIncrement;not null"`
	FollowerID uint `gorm:"not null"`
	FolloweeID uint `gorm:"not null"`
	Follower   User `gorm:"foreignKey:FollowerID"`
	Followee   User `gorm:"foreignKey:FolloweeID"`
}
