package services

import (
	"checker-api/internal/domain"
	"checker-api/internal/domain/entities"
	"checker-api/internal/ports/input/service"
	"checker-api/internal/ports/output/repository"
	"checker-api/internal/ports/output/utility"
	"context"
)

// ComparisonServiceImpl implements the ComparisonService port
// Acts as an adapter between the input layer (SQS) and the application layer (use cases)
// Difference: ComparisonService includes context.Context (application interface),
// while domain ports do not include it
type ComparisonServiceImpl struct {
	messageParser        utility.MessageParser
	responseComparator   domain.ResponseComparator
	differenceCalculator domain.DifferenceCalculator
	comparisonRepo       repository.ComparisonRepository
	idGenerator          utility.IDGenerator
}

// NewComparisonService creates a new instance of the comparison service
func NewComparisonService(
	messageParser utility.MessageParser,
	responseComparator domain.ResponseComparator,
	differenceCalculator domain.DifferenceCalculator,
	comparisonRepo repository.ComparisonRepository,
	idGenerator utility.IDGenerator,
) service.ComparisonService {
	return &ComparisonServiceImpl{
		messageParser:        messageParser,
		responseComparator:   responseComparator,
		differenceCalculator: differenceCalculator,
		comparisonRepo:       comparisonRepo,
		idGenerator:          idGenerator,
	}
}

// ProcessMessage processes a received message:
// 1. Parses and validates the message
// 2. Compares the responses
// 3. Calculates the differences
// 4. Persists the result
// TODO: Improve implementation:
// 1. Add logging for each step
// 2. Add more detailed error handling
// 3. Add metrics/observability
// 4. Consider transactions if necessary
// 5. Add additional validations
func (s *ComparisonServiceImpl) ProcessMessage(ctx context.Context, rawMessage []byte) (*entities.Comparison, error) {
	// 1. Parse and validate the message
	message, err := s.messageParser.Parse(rawMessage)
	if err != nil {
		// TODO: Return more descriptive error
		return nil, err
	}

	if err := s.messageParser.Validate(message); err != nil {
		// TODO: Return more descriptive error
		return nil, err
	}

	// 2. Compare the responses
	isEqual, err := s.responseComparator.Compare(message.Source, message.Destination)
	if err != nil {
		// TODO: Handle comparison error
		return nil, err
	}

	// 3. Calculate the differences
	differences, err := s.differenceCalculator.CalculateDifferences(message.Source, message.Destination)
	if err != nil {
		// TODO: Handle difference calculation error
		return nil, err
	}

	// 4. Create the comparison
	comparison := &entities.Comparison{
		ID:          s.idGenerator.Generate(),
		Source:      message.Source,
		Destination: message.Destination,
		Differences: differences,
		IsEqual:     isEqual,
	}

	// 5. Persist the result
	if err := s.comparisonRepo.Save(ctx, comparison); err != nil {
		// TODO: Handle persistence error (consider retry, etc.)
		return nil, err
	}

	return comparison, nil
}
