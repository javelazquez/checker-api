package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	endpointURL = "http://localhost:4566"
	region      = "us-east-1"
	maxRetries  = 30
)

func main() {
	log.Println("🚀 Inicializando recursos de LocalStack...")

	ctx := context.Background()

	// Configurar AWS SDK con endpoint personalizado para LocalStack
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		log.Fatalf("❌ Error al cargar configuración AWS: %v", err)
	}

	// Esperar a que LocalStack esté listo
	log.Println("⏳ Esperando a que LocalStack esté listo...")
	if !waitForLocalStack(ctx, cfg) {
		log.Println("⚠️  Timeout esperando a LocalStack, continuando...")
		os.Exit(0)
	}

	// Crear cliente de DynamoDB con endpoint personalizado (método no deprecado)
	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpointURL)
	})

	// Crear tabla de DynamoDB
	log.Println("📊 Creando tabla DynamoDB: comparisons")
	createDynamoTable(ctx, dynamoClient)

	// Crear cliente de SQS con endpoint personalizado (método no deprecado)
	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(endpointURL)
	})

	// Crear cola de SQS
	log.Println("📬 Creando cola SQS: comparison-queue")
	createSQSQueue(ctx, sqsClient)

	log.Println("✅ Inicialización de LocalStack completada")
}

func waitForLocalStack(ctx context.Context, cfg aws.Config) bool {
	// Crear cliente de DynamoDB con endpoint personalizado para verificar
	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpointURL)
	})
	for i := 0; i < maxRetries; i++ {
		_, err := dynamoClient.ListTables(ctx, &dynamodb.ListTablesInput{})
		if err == nil {
			log.Println("✅ LocalStack está listo")
			return true
		}
		time.Sleep(2 * time.Second)
	}
	return false
}

func createDynamoTable(ctx context.Context, client *dynamodb.Client) {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("comparisons"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}

	_, err := client.CreateTable(ctx, input)
	if err != nil {
		// Verificar si es porque la tabla ya existe
		if _, ok := err.(*types.ResourceInUseException); ok {
			log.Println("ℹ️  La tabla DynamoDB 'comparisons' ya existe")
			return
		}
		// Verificar el mensaje de error como alternativa
		errMsg := err.Error()
		if errMsg == "ResourceInUseException: Cannot create preexisting table" ||
			errMsg == "ResourceInUseException" {
			log.Println("ℹ️  La tabla DynamoDB 'comparisons' ya existe")
			return
		}
		log.Printf("⚠️  Error al crear la tabla DynamoDB: %v", err)
	} else {
		log.Println("✅ Tabla DynamoDB 'comparisons' creada exitosamente")
	}
}

func createSQSQueue(ctx context.Context, client *sqs.Client) {
	queueName := "comparison-queue"

	// Primero verificar si la cola ya existe
	getURLInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	_, err := client.GetQueueUrl(ctx, getURLInput)
	if err == nil {
		log.Printf("ℹ️  La cola SQS '%s' ya existe", queueName)
		return
	}

	// La cola no existe, crearla
	log.Printf("La cola no existe, creándola...")
	createInput := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}

	createOutput, err := client.CreateQueue(ctx, createInput)
	if err != nil {
		log.Printf("⚠️  Error al crear la cola SQS: %v", err)

		// Intentar listar todas las colas para debug
		listOutput, listErr := client.ListQueues(ctx, &sqs.ListQueuesInput{})
		if listErr == nil {
			log.Printf("Colas disponibles: %v", listOutput.QueueUrls)
		}
		return
	}

	if createOutput.QueueUrl != nil {
		log.Printf("✅ Cola SQS creada exitosamente: %s", *createOutput.QueueUrl)

		// Verificar que se creó correctamente
		time.Sleep(2 * time.Second)
		verifyOutput, verifyErr := client.GetQueueUrl(ctx, getURLInput)
		if verifyErr == nil && verifyOutput.QueueUrl != nil {
			log.Printf("✅ Cola SQS verificada: %s", *verifyOutput.QueueUrl)
		}
	} else {
		log.Println("⚠️  Cola SQS creada pero no se obtuvo la URL")
	}
}
