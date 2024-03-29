package storage

import (
	"context"
	s "strings"
	"time"

	reply "github.com/deejcoder/myprojects/reply"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Project represents a Project
type Project struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Status       string             `json:"status" bson:"status"`
	Tags         []string           `json:"tags" bson:"tags"`
	DateCreated  time.Time          `json:"date_created" bson:"date_created"`
	DateModified time.Time          `json:"date_modified" bson:"date_modified"`
	Summary      string             `json:"summary" bson:"summary"`
	ProjectLink  string             `json:"projectLink" bson:"projectLink"`
	Content      string             `json:"content" bson:"content"`
}

func setCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("projects")
}

// Validate validates if all project fields are valid, returns list of errors
func (p *Project) Validate(r *reply.Response) bool {

	if len(p.Title) < 5 || len(p.Title) > 80 {
		r.AddValidationError("title", "Title should be between 5 and 80 characters long")
	}

	if !(s.Contains(p.Status, "In progress") || s.Contains(p.Status, "Completed")) {
		r.AddValidationError("status", "The project status must be 'in progress' or 'completed'")
	}

	if len(p.Tags) > 8 {
		r.AddValidationError("tags", "You may only have 8 project tags")
	}

	if len(p.Summary) < 50 || len(p.Summary) > 350 {
		r.AddValidationError("summary", "Project summary should be between 50 and 350 characters long")
	}

	if len(p.Content) < 50 || len(p.Content) > 10000 {
		r.AddValidationError("content", "The project content should be between 50 and 10000 characters")
	}

	return !r.HasValidationErrors()
}

// GetProject finds a project by project ID, and returns it. If not found, nil.
func GetProject(db *mongo.Database, id string) *Project {
	col := setCollection(db)

	oid, invalid := primitive.ObjectIDFromHex(id)
	// check if id is a valid id
	if invalid != nil {
		log.Error(invalid)
		return nil
	}

	filter := bson.M{"_id": oid}
	var project *Project
	if err := col.FindOne(context.TODO(), filter).Decode(&project); err != nil {
		log.Error(err)
		return nil
	}
	return project
}

// GetProjects returns all projects
func GetProjects(db *mongo.Database) []*Project {
	col := setCollection(db)

	projects := make([]*Project, 0)
	cursor, err := col.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Error(err)
		return projects
	}

	defer cursor.Close(context.TODO())

	// fetch a project from mongo
	for cursor.Next(context.TODO()) {
		var project Project
		if err := cursor.Decode(&project); err != nil {
			log.Error(err)
			continue
		}

		projects = append(projects, &project)
	}

	if err := cursor.Err(); err != nil {
		log.Error(err)
	}
	return projects
}

// UpdateProject updates an existing project, returns a list of form errors & if project was updated
func UpdateProject(db *mongo.Database, project *Project) bool {
	col := setCollection(db)

	filter := bson.M{"_id": project.ID}
	result := col.FindOneAndReplace(context.TODO(), filter, project)
	if result.Err() != nil {
		log.Error(result.Err())
		return false
	}
	return true
}

// DeleteProject removes an existing project
func DeleteProject(db *mongo.Database, id string) bool {
	col := setCollection(db)

	pid, invalid := primitive.ObjectIDFromHex(id)
	if invalid != nil {
		return false
	}

	filter := bson.M{"_id": pid}

	doc := col.FindOneAndDelete(context.TODO(), filter)
	if err := doc.Err(); err != nil {
		log.Error(err)
		return false
	}
	return true
}
