package controllers

import (
	"fmt"
	"net/http"
	"context"
	"encoding/json"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
	
	"../database"
	"../models"
	"../util"
	"../middlewares"
)

//GetPlacesByUserID .
func GetPlacesByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := mux.Vars(r)["uid"]
	var places models.Places
	ObjctID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		util.Err(w, err, http.StatusBadRequest)
		return
	} 
	cur, err := database.PlacesColl.Find(context.Background(), bson.D{{Key:"creator",Value: ObjctID}})
	if err != nil {
    util.Err(w, err, http.StatusNotFound)
		return
	}
	if err = cur.All(context.TODO(), &places.Places); err != nil {
		util.Err(w, err, http.StatusInternalServerError)
		return 
	}
	json.NewEncoder(w).Encode(places)
}

//GetPlaceByID .
func GetPlaceByID (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pid := mux.Vars(r)["pid"]
	var place models.Place
	ObjctID, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		util.Err(w, err, http.StatusBadRequest)
		return
	} 
	result := database.PlacesColl.FindOne(context.Background(), bson.D{{Key:"_id",Value: ObjctID}})
	if err := result.Decode(&place); err != nil {
		util.Err(w, err, http.StatusNotFound)
		return 
	}
	json.NewEncoder(w).Encode(place)
}


//CreatePlace .
func CreatePlace (w http.ResponseWriter, r * http.Request) {
	//...後でvalidationを追加する
	w.Header().Set("Content-Type", "application/json")
	
	var place models.Place
	json.NewDecoder(r.Body).Decode(&place)

	//.... convert address to location
	location, err := util.GetCoordsForAddress(place.Address)
	if err != nil {
			util.Err(w, err, http.StatusBadRequest)
			return
	}
	place.Location = location
	
	//jwtからuserIdを取得
	UserObjctID, err := primitive.ObjectIDFromHex(middlewares.UserID) //仮のuserId 実際はjwtからもってくる
	if err != nil {
			util.Err(w, err, http.StatusBadRequest)
			return
	}


	//transaction place追加 and user.placesにplace_idを追加
	err = database.Client.UseSession(context.Background(), func (sessCtx mongo.SessionContext) error {

		//transaction start!
		if err := sessCtx.StartTransaction(); err != nil {
			return err
		}
		
		result, err := database.PlacesColl.InsertOne(sessCtx, place)
		if err != nil {
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		var user models.User
		update := bson.D{{Key:"$push",Value: bson.D{{Key:"places",Value: result.InsertedID}}}}
		err = database.UserColl.FindOneAndUpdate(sessCtx, bson.D{{Key:"_id",Value: UserObjctID}}, update).Decode(&user)
		if err != nil {
			fmt.Println(err)
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		err = sessCtx.CommitTransaction(sessCtx)
		return err
	})
	if err != nil {
		util.Err(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(place)

	//... read creator from jwt
	//... image adding
	//返す値の確認
}


//UpdatePlace .
func UpdatePlace (w http.ResponseWriter, r *http.Request) {
	//後でvalidation追加
	w.Header().Set("Content-Type", "application/json")

	pid := mux.Vars(r)["pid"]
	PlaceObjctID, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		util.Err(w, err, http.StatusBadRequest)
		return
	}

	var updatePlace struct {
		Title       string
		Description string 
	}
	json.NewDecoder(r.Body).Decode(&updatePlace)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: updatePlace.Title}, {Key: "description", Value: updatePlace.Description}}}}
	var place models.Place

	err = database.PlacesColl.FindOneAndUpdate(context.Background(), bson.D{{Key: "_id", Value: PlaceObjctID}}, update).Decode(&place)
	if err != nil {
		util.Err(w, err, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(place)
}


//DeletePlace .
func DeletePlace (w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

	pid := mux.Vars(r)["pid"]
	PlaceObjctID, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		util.Err(w, err, http.StatusBadRequest)
		return
	}
	var deletePlace models.Place
	var user models.User
	
	//transaction
	err = database.Client.UseSession(context.Background(), func (sessCtx mongo.SessionContext) error {
		//transaction start!
		if err := sessCtx.StartTransaction(); err != nil {
			return err
		}
		err := database.PlacesColl.FindOneAndDelete(sessCtx, bson.D{{Key: "_id", Value: PlaceObjctID}}).Decode(&deletePlace)
		if err != nil {
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		userID := deletePlace.Creator
		update := bson.D{{Key:"$pull",Value: bson.D{{Key: "places", Value: PlaceObjctID}}}}//test
		err = database.UserColl.FindOneAndUpdate(sessCtx, bson.D{{Key: "_id", Value: userID}}, update).Decode(&user)
		if err != nil {
			fmt.Println(err)
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		err = sessCtx.CommitTransaction(sessCtx)
		return err
	})

	if err != nil {
		util.Err(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deletePlace)
}