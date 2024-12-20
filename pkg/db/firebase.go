package db

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type FirebaseDatabaseConfig struct {
	Type                    string `json:"type" mapstructure:"type"`
	ProjectId               string `json:"project_id" mapstructure:"project_id"`
	PrivateKeyId            string `json:"private_key_id" mapstructure:"private_key_id"`
	PrivateKey              string `json:"private_key" mapstructure:"private_key"`
	ClientEmail             string `json:"client_email" mapstructure:"client_email"`
	ClientId                string `json:"client_id" mapstructure:"client_id"`
	AuthUri                 string `json:"auth_uri" mapstructure:"auth_uri"`
	TokenUri                string `json:"token_uri" mapstructure:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url" mapstructure:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url" mapstructure:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain" mapstructure:"universe_domain"`
	StorageBucket           string `json:"storage_bucket" mapstructure:"storage_bucket"`
}

type FirebaseDatabaseClient struct {
	Client *firestore.Client
}

type FirebaseDatabaseInterface interface {
	Close()
	Collection(name string) *firestore.CollectionRef
	Documents(ctx context.Context, name string) *firestore.DocumentIterator
}

func NewFirebaseDatabaseConnection(ctx context.Context, fields any, kind string) (FirebaseDatabaseInterface, error) {

	fbConfig, err := parseConfig(fields)
	if err != nil {
		return nil, err
	}

	// // // Converta o mapa para JSON
	dbJSON, err := json.Marshal(fbConfig)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(dbJSON)
	_, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := firestore.NewClient(ctx, fbConfig.ProjectId, opt)
	if err != nil {
		log.Fatalf("Error creating Firestore client: %v\n", err)
	}

	query := FirebaseDatabaseClient{
		Client: client,
	}

	return &query, nil
}

func (fbd *FirebaseDatabaseClient) Close() {
	fbd.Client.Close()
}

func (fbd *FirebaseDatabaseClient) Collection(name string) *firestore.CollectionRef {
	return fbd.Client.Collection(name)
}

func (fbd *FirebaseDatabaseClient) Documents(ctx context.Context, name string) *firestore.DocumentIterator {
	return fbd.Client.Collection(name).Documents(ctx)
}

func parseConfig(fields any) (*FirebaseDatabaseConfig, error) {

	b, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	var config FirebaseDatabaseConfig
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
