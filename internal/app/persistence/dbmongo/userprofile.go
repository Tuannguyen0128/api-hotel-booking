package dbmongo

import (
	"api-hotel-booking/internal/app/persistence"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u *UserRepo) Insert(ctx context.Context, profile persistence.UserProfile) (string, error) {
	_, err := u.GetActiveInActiveUserProfilesByEmail(ctx, profile.Email)
	if err == nil {
		return "", persistence.DuplicateEmailError
	} else {
		if err != persistence.NotFoundError {
			return "", err
		}
	}
	profile.Id = "user_" + uuid.NewString()
	if err := u.client.Insert(ctx, dbName, collUserProfile, profile); err != nil {
		return "", castKnownError(err)
	}
	return profile.Id, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (persistence.UserProfile, error) {
	profile := persistence.UserProfile{}
	query := bson.D{
		{Key: "email", Value: email},
	}
	collection := u.client.GetCollection(dbName, collUserProfile)
	err := collection.FindOne(ctx, query, &options.FindOneOptions{Sort: bson.D{{"createDt", -1}}}).Decode(&profile)
	if err != nil {
		return profile, castKnownError(err)
	}
	return profile, nil
}

func (u *UserRepo) GetById(ctx context.Context, userId string) (persistence.UserProfile, error) {
	profile := persistence.UserProfile{}

	err := u.client.SearchOne(ctx, dbName, collUserProfile, bson.D{{"_id", userId}}, &profile)
	if err != nil {
		return profile, castKnownError(err)
	}
	return profile, nil
}

func (u *UserRepo) GetByCompanyId(ctx context.Context, companyId string, filter persistence.UserFilter) ([]persistence.UserProfile, error) {
	var profiles []persistence.UserProfile
	query := make(bson.M)
	if companyId != "" {
		query["companyId"] = companyId
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	err := u.client.Search(ctx, dbName, collUserProfile, query, &profiles)

	if err != nil {
		return profiles, castKnownError(err)
	}
	return profiles, nil
}

func (u *UserRepo) GetAll(ctx context.Context, filter persistence.UserFilter) ([]persistence.UserProfile, error) {
	var profiles []persistence.UserProfile
	query := make(bson.M)
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	err := u.client.Search(ctx, dbName, collUserProfile, query, &profiles)

	if err != nil {
		return profiles, castKnownError(err)
	}
	return profiles, nil
}

func (u *UserRepo) Update(ctx context.Context, userId string, document persistence.EditUserProfile) error {
	query := bson.D{{"_id", userId}}
	update := bson.D{{"$set", document}}
	err := u.client.Update(ctx, dbName, collUserProfile, query, update)
	if err != nil {
		return castKnownError(err)
	}
	return nil
}

func (u *UserRepo) UpdateInfo(ctx context.Context, updateReq persistence.EditUserProfileInfo) error {
	query := bson.D{{Key: "_id", Value: updateReq.UserId}, {Key: "companyId", Value: updateReq.CompanyId}}
	if updateReq.Role != "" {
		query = append(query, bson.E{Key: "role", Value: updateReq.Role})
	}

	update := bson.D{{"$set", updateReq}}

	if err := u.client.Update(ctx, dbName, collUserProfile, query, update); err != nil {
		return castKnownError(err)
	}
	return nil
}

func (u *UserRepo) WrongPasswordCounterIncrease(userId string, lockIfOver int, changeStatusTo string) (int, error) {
	query := bson.D{{Key: "_id", Value: userId}}

	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "wrong_password", Value: 1},
		}},
	}

	projection := bson.D{{Key: "wrong_password", Value: 1}, {Key: "companyId", Value: 1}}

	var result persistence.UserProfile
	if err := u.client.GetCollection(dbName, collUserProfile).FindOneAndUpdate(context.TODO(), query, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After), options.FindOneAndUpdate().SetProjection(projection)).Decode(&result); err != nil {
		return 0, castKnownError(err)
	}

	if result.WrongPassword >= lockIfOver {
		return result.WrongPassword, u.UpdateInfo(context.TODO(), persistence.EditUserProfileInfo{
			UserId:    userId,
			CompanyId: result.CompanyId,
			NewStatus: &changeStatusTo,
		})
	}

	return result.WrongPassword, nil
}

func (u *UserRepo) WrongPasswordCounterReset(userId string) error {
	query := bson.D{{Key: "_id", Value: userId}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "wrong_password", Value: 0},
	}}}

	if result, err := u.client.GetCollection(dbName, collUserProfile).UpdateOne(context.TODO(), query, update); err != nil {
		return castKnownError(err)
	} else if result.ModifiedCount == 0 {
		return persistence.NotFoundError
	}

	return nil
}

func (u *UserRepo) GetActiveInActiveUserProfilesByEmail(ctx context.Context, email string) ([]persistence.UserProfile, error) {
	profiles := []persistence.UserProfile{}
	query := bson.D{
		{Key: "email", Value: email},
		{Key: "status", Value: bson.D{{Key: "$ne", Value: user_delete}}},
	}
	err := u.client.Search(ctx, dbName, collUserProfile, query, &profiles)
	if err != nil {
		return profiles, castKnownError(err)
	}
	return profiles, nil
}
