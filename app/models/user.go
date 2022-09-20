package models

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/storage"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	usernameReg = "^[A-Za-z\\d]\\w{1,19}$"
	emailReg    = "^\\w+@[a-z0-9]+(\\.[a-z]+){1,3}$"
)

type User struct {
	UID            uint64             `bson:"_id"`
	Username       string             `bson:"username,omitempty"`
	Misesid        string             `bson:"misesid,omitempty"`
	Gender         enum.Gender        `bson:"gender,misesid"`
	Mobile         string             `bson:"mobile,omitempty"`
	Email          string             `bson:"email,omitempty"`
	Address        string             `bson:"address,omitempty"`
	Intro          string             `bson:"intro,omitempty"`
	AvatarPath     string             `bson:"avatar_path,omitempty"`
	FollowingCount uint32             `bson:"following_count,omitempty"`
	FansCount      uint32             `bson:"fans_count,omitempty"`
	LikedCount     uint32             `bson:"liked_count,omitempty"`
	CreatedAt      time.Time          `bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty"`
	OnChain        bool               `bson:"on_chain,omitempty"`
	ChannelID      primitive.ObjectID `bson:"channel_id,omitempty"`
	NftAvatar      *NftAvatar         `bson:"nft_avatar,omitempty"`
	AvatarUrl      string             `bson:"-"`
	IsFollowed     bool               `bson:"-"`
	IsAirdropped   bool               `bson:"-"`
	IsLogined      bool               `bson:"-"`
	AirdropStatus  bool               `bson:"-"`
	IsFriend       bool               `bson:"-"`
	Tags           []enum.TagType     `bson:"tags"`
	IsBlocked      bool               `bson:"-"`
	NewFansCount   uint32             `bson:"-"`
	RelationType   enum.RelationType  `bson:"-"`
	BlockState     enum.BlockState    `bson:"-"`
	Avatar         *Avatar            `bson:"-"`
	Pubkey         string             `bson:"pubkey,omitempty"`
	EthAddress     string             `bson:"eth_address,omitempty"`
}

type NftAvatar struct {
	NftAssetID        primitive.ObjectID `bson:"nft_asset_id"`
	ImageURL          string             `bson:"image_url"`
	ImagePreviewUrl   string             `bson:"image_preview_url"`
	ImageThumbnailUrl string             `bson:"image_thumbnail_url"`
}
type Avatar struct {
	Orgin  string
	Large  string
	Medium string
	Small  string
}

func (u *User) Validate(ctx context.Context) error {
	if err := u.validateUsername(ctx); err != nil {
		return err
	}
	if err := u.validateEmail(ctx); err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeCreate(ctx context.Context) error {
	var err error
	u.UID, err = getNextSeq(ctx, "userid")
	if err != nil {
		return err
	}
	u.CreatedAt = time.Now()
	return u.BeforeUpdate(ctx)
}

func (u *User) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	if err := u.Validate(ctx); err != nil {
		return err
	}
	return nil
}

func (u *User) UpdatePostTime(ctx context.Context, t time.Time) error {
	return db.DB().Collection("users").FindOneAndUpdate(ctx, bson.M{"_id": u.UID},
		bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "latest_post_time",
				Value: &t,
			}, {
				Key:   "updated_at",
				Value: time.Now(),
			}}},
		}).Err()
}

func ListUserByIDs(ctx context.Context, uids ...uint64) ([]*User, error) {
	users := make([]*User, 0)
	chain := db.ODM(ctx).Where(bson.M{
		"_id": bson.M{"$in": uids},
	})
	return users, chain.Find(&users).Error
}

func FindUser(ctx context.Context, uid uint64) (*User, error) {
	user := &User{}
	result := db.DB().Collection("users").FindOne(ctx, &bson.M{
		"_id": uid,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	return user, result.Decode(user)
}

func FindOrCreateUserByMisesid(ctx context.Context, misesid, pubkey string) (*User, bool, error) {
	user := &User{}
	result := db.DB().Collection("users").FindOne(ctx, &bson.M{
		"misesid": misesid,
	})
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		created, err := createMisesUser(ctx, misesid, pubkey)
		return created, true, err
	}
	if err != nil {
		return nil, false, err
	}
	return user, false, result.Decode(user)
}

func UpdateUserProfile(ctx context.Context, user *User) error {
	err := user.BeforeUpdate(ctx)
	if err != nil {
		return err
	}
	_, err = db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"gender":     user.Gender,
			"mobile":     user.Mobile,
			"email":      user.Email,
			"address":    user.Address,
			"intro":      user.Intro,
			"updated_at": time.Now(),
		}}})
	return err
}
func UpdateUserEthAdress(ctx context.Context, user *User) error {
	err := user.BeforeUpdate(ctx)
	if err != nil {
		return err
	}
	_, err = db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"pubkey":      user.Pubkey,
			"eth_address": user.EthAddress,
			"updated_at":  time.Now(),
		}}})
	return err
}

func UpdateUsername(ctx context.Context, user *User) error {
	err := user.BeforeUpdate(ctx)
	if err != nil {
		return err
	}
	_, err = db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"username":   user.Username,
			"updated_at": time.Now(),
		}}})
	return err
}

func UpdateUserAvatar(ctx context.Context, user *User) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"avatar_path": user.AvatarPath,
			"nft_avatar":  nil,
			"updated_at":  time.Now(),
		}}})
	return err
}
func UpdateUserNftAvatar(ctx context.Context, user *User) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"nft_avatar": user.NftAvatar,
			"updated_at": time.Now(),
		}}})
	return err
}
func UpdateUserOnChainByMisesid(ctx context.Context, misesid string) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"misesid": misesid,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"on_chain": true,
		}}})
	return err
}
func UpdateUserChannelIDByUID(ctx context.Context, uid uint64, channel_id primitive.ObjectID) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": uid,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"channel_id": channel_id,
		}}})
	return err
}

func createMisesUser(ctx context.Context, misesid, pubkey string) (*User, error) {
	user := &User{
		Misesid: misesid,
		Pubkey:  pubkey,
	}
	err := user.BeforeCreate(ctx)
	if err != nil {
		return nil, err
	}
	address, err := PubkeyToEthAddress(pubkey)
	if err == nil {
		user.EthAddress = address
	}
	_, err = db.DB().Collection("users").InsertOne(ctx, user)
	return user, err
}
func FindUserEthAddress(ctx context.Context, uid uint64) (*User, error) {
	user, err := FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user.EthAddress == "" {
		chainUser, err := FindChainUser(ctx, &search.ChainUserSearch{Misesid: user.Misesid})
		if err != nil {
			return nil, err
		}
		address, err := PubkeyToEthAddress(chainUser.Pubkey)
		if err != nil {
			return nil, err
		}
		user.EthAddress = address
		user.Pubkey = chainUser.Pubkey
	}
	return user, UpdateUserEthAdress(ctx, user)
}

func PubkeyToEthAddress(pubkey string) (string, error) {
	r, err := hex.DecodeString(pubkey)
	if err != nil {
		return "", err
	}
	btcec_pubKey, err := btcec.ParsePubKey(r, btcec.S256())
	if err != nil {
		return "", err
	}
	a := btcec_pubKey.ToECDSA()
	addr := crypto.PubkeyToAddress(*a)
	return addr.Hex(), nil
}

func FindUserByIDs(ctx context.Context, ids ...uint64) ([]*User, error) {
	users := make([]*User, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": ids}}).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, PreloadUserData(ctx, users...)
}
func FindUserByMisesids(ctx context.Context, misesids ...string) ([]*User, error) {
	users := make([]*User, 0)
	err := db.ODM(ctx).Where(bson.M{"misesid": bson.M{"$in": misesids}}).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func FindUserByMisesid(ctx context.Context, misesid string) (*User, error) {
	user := &User{}
	err := db.ODM(ctx).Where(bson.M{"misesid": utils.AddMisesidProfix(misesid)}).Last(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func FindUserByEthAddress(ctx context.Context, address string) (*User, error) {
	user := &User{}
	err := db.ODM(ctx).Where(bson.M{"eth_address": utils.EthAddressToEIPAddress(address)}).Last(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func FindUserByEthAddresses(ctx context.Context, addresses ...string) ([]*User, error) {
	for k, v := range addresses {
		addresses[k] = utils.EthAddressToEIPAddress(v)
	}
	fmt.Println(addresses)
	res := make([]*User, 0)
	err := db.ODM(ctx).Where(bson.M{"eth_address": bson.M{"$in": addresses}}).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func GetUserMapByEthAddresses(ctx context.Context, addresses ...string) (map[string]*User, error) {
	users, err := FindUserByEthAddresses(ctx, addresses...)
	if err != nil {
		return nil, err
	}
	addressMap := make(map[string]*User)
	for _, user := range users {
		addressMap[strings.ToLower(user.EthAddress)] = user
	}
	return addressMap, nil
}

func GetUserMap(ctx context.Context, ids ...uint64) (map[uint64]*User, error) {
	users, err := FindUserByIDs(ctx, ids...)
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]*User)
	for _, user := range users {
		result[user.UID] = user
	}
	return result, nil
}

func PreloadUserData(ctx context.Context, users ...*User) error {
	err := PreloadUserAvatar(ctx, users...)
	if err != nil {
		return err
	}
	err = preloadCurrentUserRelationship(ctx, users...)
	if err != nil {
		return err
	}
	return nil
}

func PreloadUserAvatar(ctx context.Context, users ...*User) error {
	paths := make([]string, 0)
	for _, user := range users {
		if user.AvatarPath != "" {
			paths = append(paths, user.AvatarPath)
		}
	}
	avatars, err := storage.ImageClient.GetFileUrl(ctx, paths...)
	if err != nil {
		return err
	}
	//thumb image
	optsThumb := &storage.ImageOptions{
		ResizeOptions: &storage.ResizeOptions{Resize: true, Width: 128, Height: 128},
	}
	optsMedium := &storage.ImageOptions{
		ResizeOptions: &storage.ResizeOptions{Resize: true, Width: 200, Height: 200},
	}
	thumbImages, err := storage.ImageClient.GetFileUrlOptions(ctx, optsThumb, paths...)
	mediumImages, err := storage.ImageClient.GetFileUrlOptions(ctx, optsMedium, paths...)
	for _, user := range users {
		user.AvatarUrl = avatars[user.AvatarPath]
		user.Avatar = &Avatar{}
		user.Avatar.Orgin = avatars[user.AvatarPath]
		user.Avatar.Large = avatars[user.AvatarPath]
		user.Avatar.Small = thumbImages[user.AvatarPath]
		user.Avatar.Medium = mediumImages[user.AvatarPath]
	}
	return nil
}

func preloadCurrentUserRelationship(ctx context.Context, users ...*User) error {
	currentUID := ctx.Value(utils.CurrentUIDKey{})
	if currentUID == nil {
		return nil
	}
	uid := currentUID.(uint64)
	if uid == 0 {
		return nil
	}
	toUIDs := make([]uint64, len(users))
	for i, user := range users {
		toUIDs[i] = user.UID
	}
	followMap, err := GetFollowMap(ctx, uid, toUIDs)
	if err != nil {
		return err
	}
	blacklistMap, err := GetBlacklistMap(ctx, uid, toUIDs)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.IsFollowed = followMap[user.UID] != nil
		user.IsFriend = followMap[user.UID] != nil && followMap[user.UID].IsFriend
		user.IsBlocked = blacklistMap[user.UID] != nil
	}
	return nil
}

func (u *User) validateUsername(ctx context.Context) error {
	if u.Username == "" {
		return nil
	}
	match, _ := regexp.MatchString(usernameReg, u.Username)
	if !match {
		return codes.ErrUnprocessableEntity
	}
	query := db.ODM(ctx).Where(bson.M{"username": u.Username})
	if u.UID != 0 {
		query = query.Where(bson.M{"_id": bson.M{"$ne": u.UID}})
	}
	var c int64
	err := query.Model(u).Count(&c).Error
	if err != nil {
		return err
	}
	if c > 0 {
		return codes.ErrUsernameDuplicate
	}
	return nil
}

func (u *User) validateEmail(ctx context.Context) error {
	if u.Email == "" {
		return nil
	}
	match, _ := regexp.MatchString(emailReg, u.Email)
	if !match {
		return codes.ErrUnprocessableEntity.New("invalid email")
	}
	return nil
}
