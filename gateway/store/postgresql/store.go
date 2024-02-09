package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/gateway/store/clouddatastore/models"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
)

// Implements the Store interface
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new Store
func NewStore(ctx context.Context) (*Store, error) {
	return &Store{db: nil}, nil
}

// FirstEvent implements store.Store.
func (s *Store) FirstEvent(ctx context.Context) (*types.GatewayEvent, error) {
	row := s.db.QueryRowxContext(ctx, "SELECT * FROM gateway_events ORDER BY block_number, transaction_index, log_index LIMIT 1")
	var e models.DBGatewayEvent
	if err := row.StructScan(&e); err != nil {
		return nil, err
	}
	return e.GatewayEvent(), nil
}

// CleanOldPendingEvents implements store.Store.
func (s *Store) CleanOldPendingEvents(ctx context.Context, height uint64) error {
	res := s.db.MustExecContext(ctx, "DELETE FROM pending_gateway_events WHERE block_number < $1", height)
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// CurrentBlock implements store.Store.
func (s *Store) CurrentBlock(ctx context.Context, process string) (uint64, error) {
	row := s.db.QueryRowxContext(ctx, "SELECT height FROM current_blocks WHERE process = $1", process)
	var height uint64
	if err := row.Scan(&height); err != nil {
		return 0, err
	}
	return height, nil
}

// StoreCurrentBlock implements store.Store.
func (s *Store) StoreCurrentBlock(ctx context.Context, process string, height uint64) error {
	res := s.db.MustExecContext(ctx, "INSERT INTO current_blocks (process, height) VALUES ($1, $2) ON CONFLICT (process) DO UPDATE SET height = $2", process, height)
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// Delete implements store.Store.
func (s *Store) Delete(ctx context.Context, id types.ID) error {
	res := s.db.MustExecContext(ctx, "DELETE FROM gateways WHERE id = $1", id.String())
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// DeletePendingEvent implements store.Store.
func (s *Store) DeletePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error {
	res := s.db.MustExecContext(ctx, "DELETE FROM pending_gateway_events WHERE block_number = $1 AND transaction_index = $2 AND log_index = $3", pendingEvent.BlockNumber, pendingEvent.TransactionIndex, pendingEvent.LogIndex)
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// EventsFromTo implements store.Store.
func (s *Store) EventsFromTo(ctx context.Context, from uint64, to uint64) ([]*types.GatewayEvent, error) {
	query := "SELECT * FROM gateway_events WHERE block_number >= $1 AND block_number <= $2"
	rows, err := s.db.QueryxContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*types.GatewayEvent
	for rows.Next() {
		var e models.DBGatewayEvent
		if err := rows.StructScan(&e); err != nil {
			return nil, err
		}
		events = append(events, e.GatewayEvent())
	}
	return events, nil
}

// Get implements store.Store.
func (s *Store) Get(ctx context.Context, id types.ID) (*types.Gateway, error) {
	query := "SELECT * FROM gateways WHERE id = $1"
	row := s.db.QueryRowxContext(ctx, query, id.String())
	var dbGateway models.DBGateway
	if err := row.StructScan(&dbGateway); err != nil {
		return nil, err
	}
	return dbGateway.Gateway(), nil
}

// GetAll implements store.Store.
func (s *Store) GetAll(ctx context.Context) ([]*types.Gateway, error) {
	query := "SELECT * FROM gateways"
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gateways []*types.Gateway
	for rows.Next() {
		var dbGateway models.DBGateway
		if err := rows.StructScan(&dbGateway); err != nil {
			return nil, err
		}
		gateways = append(gateways, dbGateway.Gateway())
	}
	return gateways, nil
}

// GetByOwner implements store.Store.
func (s *Store) GetByOwner(ctx context.Context, owner common.Address, limit int, cursor string) ([]*types.Gateway, string, error) {
	query := "SELECT * FROM gateways WHERE owner = $1"
	args := []interface{}{owner.String()}

	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	if cursor != "" {
		query += " AND id > $3"
		args = append(args, cursor)
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var gateways []*types.Gateway
	for rows.Next() {
		var dbGateway models.DBGateway
		if err := rows.StructScan(&dbGateway); err != nil {
			return nil, "", err
		}
		gateways = append(gateways, dbGateway.Gateway())
	}

	// Get the last gateway ID as the cursor
	lastGatewayID := ""
	if len(gateways) > 0 {
		lastGatewayID = gateways[len(gateways)-1].ID.String()
	}

	return gateways, lastGatewayID, nil
}

// GetCountInCellAtRes implements store.Store.
func (*Store) GetCountInCellAtRes(ctx context.Context, cell h3light.Cell, res int) (map[h3light.Cell]uint64, error) {

}

// GetEvents implements store.Store.
func (s *Store) GetEvents(ctx context.Context, gatewayID types.ID, limit int, cursor string) ([]*types.GatewayEvent, string, error) {
	query := "SELECT * FROM gateway_events WHERE gateway_id = $1"
	args := []interface{}{gatewayID.String()}
	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}
	if cursor != "" {
		query += " AND id > $3"
		args = append(args, cursor)
	}
	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var events []*types.GatewayEvent
	for rows.Next() {
		var e models.DBGatewayEvent
		if err := rows.StructScan(&e); err != nil {
			return nil, "", err
		}
		events = append(events, e.GatewayEvent())
	}

	// Get the last event ID as the cursor
	lastEventID := ""
	if len(events) > 0 {
		lastEventID = events[len(events)-1].ID.String()
	}

	return events, lastEventID, nil
}

// GetEventsBetween implements store.Store.
func (s *Store) GetEventsBetween(ctx context.Context, start time.Time, end time.Time) ([]*types.GatewayEvent, error) {
	query := "SELECT * FROM gateway_events WHERE time >= $1 AND time <= $2"
	rows, err := s.db.QueryxContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*types.GatewayEvent
	for rows.Next() {
		var e models.DBGatewayEvent
		if err := rows.StructScan(&e); err != nil {
			return nil, err
		}
		events = append(events, e.GatewayEvent())
	}
	return events, nil
}

// GetGatewayOnboardByGatewayID implements store.Store.
func (s *Store) GetGatewayOnboardByGatewayID(ctx context.Context, gatewayID string) (*models.GatewayOnboard, error) {
	query := "SELECT * FROM gateway_onboards WHERE gateway_id = $1"
	row := s.db.QueryRowxContext(ctx, query, gatewayID)
	var onboard models.DBGatewayOnboard
	if err := row.StructScan(&onboard); err != nil {
		return nil, err
	}
	return onboard.GatewayOnboard(), nil
}

// GetGatewayOnboardsByOwner implements store.Store.
func (s *Store) GetGatewayOnboardsByOwner(ctx context.Context, onboarder common.Address, owner common.Address, limit int, cursor string) ([]*models.GatewayOnboard, string, error) {
	query := "SELECT * FROM gateway_onboards WHERE onboarder = $1 AND owner = $2"
	args := []interface{}{onboarder.String(), owner.String()}
	if limit > 0 {
		query += " LIMIT $3"
		args = append(args, limit)
	}
	if cursor != "" {
		query += " AND gateway_id > $4"
		args = append(args, cursor)
	}
	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var onboards []*models.GatewayOnboard
	for rows.Next() {
		var onboard models.DBGatewayOnboard
		if err := rows.StructScan(&onboard); err != nil {
			return nil, "", err
		}
		onboards = append(onboards, onboard.GatewayOnboard())
	}

	// Get the last gateway ID as the cursor
	lastGatewayID := ""
	if len(onboards) > 0 {
		lastGatewayID = onboards[len(onboards)-1].GatewayID
	}

	return onboards, lastGatewayID, nil
}

// GetHistoryAt implements store.Store.
func (s *Store) GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.GatewayHistory, error) {
	query := "SELECT * FROM gateway_history WHERE id = $1 AND time <= $2 ORDER BY time DESC LIMIT 1"
	row := s.db.QueryRowxContext(ctx, query, id.String(), at)
	var history models.DBGatewayHistory
	if err := row.StructScan(&history); err != nil {
		return nil, err
	}
	// Convert DBGatewayHistory to types.GatewayHistory and return it
	return history.GatewayHistory(), nil
}

// GetInCell implements store.Store.
func (s *Store) GetInCell(ctx context.Context, cell h3light.Cell) ([]*types.Gateway, error) {
	query := "SELECT * FROM gateways WHERE location LIKE $1"
	rows, err := s.db.QueryxContext(ctx, query, fmt.Sprintf("%s%%", string(cell.DatabaseCell())))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gateways []*types.Gateway
	for rows.Next() {
		var dbGateway models.DBGateway
		if err := rows.StructScan(&dbGateway); err != nil {
			return nil, err
		}
		gateways = append(gateways, dbGateway.Gateway())
	}
	return gateways, nil
}

// GetRes3CountPerRes0 implements store.Store.
func (*Store) GetRes3CountPerRes0(ctx context.Context) (map[h3light.Cell]map[h3light.Cell]uint64, error) {
	panic("unimplemented")
}

// PendingEventsForOwner implements store.Store.
func (*Store) PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.GatewayEvent, error) {
	query := "SELECT * FROM pending_gateway_events WHERE new_owner = $1 OR old_owner = $1"
	rows, err := s.db.QueryxContext(ctx, query, owner.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*types.GatewayEvent
	for rows.Next() {
		var e models.DBGatewayEvent
		if err := rows.StructScan(&e); err != nil {
			return nil, err
		}
		events = append(events, e.GatewayEvent())
	}
	return events, nil
}

// PurgeExpiredOnboards implements store.Store.
func (*Store) PurgeExpiredOnboards(ctx context.Context, expiry time.Duration) error {
	res := s.db.MustExecContext(ctx, "DELETE FROM gateway_onboards WHERE created_at < $1", time.Now().Add(-expiry))
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// Store implements store.Store.
func (s *Store) Store(ctx context.Context, gateway *types.Gateway) error {
	dbGateway := models.NewDBGateway(gateway)
	_, err := s.db.NamedExecContext(ctx, "INSERT INTO gateways (id, owner, contract_address, version, location, altitude, frequency_plan, antenna_gain, created_at, updated_at) VALUES (:id, :owner, :contract_address, :version, :location, :altitude, :frequency_plan, :antenna_gain, :created_at, :updated_at)", dbGateway)
	return err
}

// StoreEvent implements store.Store.
func (s *Store) StoreEvent(ctx context.Context, event *types.GatewayEvent) error {
	dbEvent := models.NewDBGatewayEvent(event)
	_, err := s.db.NamedExecContext(ctx, "INSERT INTO gateway_events (contract_address, block_number, transaction_index, log_index, block, transaction, type, id, version, new_owner, old_owner, new_location, old_location, new_altitude, old_altitude, new_frequency_plan, old_frequency_plan, new_antenna_gain, old_antenna_gain, time) VALUES (:contract_address, :block_number, :transaction_index, :log_index, :block, :transaction, :type, :id, :version, :new_owner, :old_owner, :new_location, :old_location, :new_altitude, :old_altitude, :new_frequency_plan, :old_frequency_plan, :new_antenna_gain, :old_antenna_gain, :time)", dbEvent)
	return err
}

// StoreGatewayOnboard implements store.Store.
func (s *Store) StoreGatewayOnboard(ctx context.Context, onboarder common.Address, gatewayID types.ID, owner common.Address, signature string, version uint8, localId string) error {
	dbOnboard := models.NewDBGatewayOnboard(gatewayID, owner, signature, version, localId, onboarder, time.Now())
	_, err := s.db.NamedExecContext(ctx, "INSERT INTO gateway_onboards (gateway_id, owner, signature, version, local_id, onboarder, created_at) VALUES (:gateway_id, :owner, :signature, :version, :local_id, :onboarder, :created_at)", dbOnboard)
	return err
}

// StoreHistory implements store.Store.
func (s *Store) StoreHistory(ctx context.Context, history *types.GatewayHistory) error {
	dbGatewayHistory := models.NewDBGatewayHistory(history)
	_, err := s.db.NamedExecContext(ctx, "INSERT INTO gateway_history (id, contract_address, version, owner, antenna_gain, frequency_plan, location, altitude, time, block_number, block, transaction) VALUES (:id, :contract_address, :version, :owner, :antenna_gain, :frequency_plan, :location, :altitude, :time, :block_number, :block, :transaction)", dbGatewayHistory)
	return err
}

// StorePendingEvent implements store.Store.
func (s *Store) StorePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error {
	dbPendingGatewayEvent := models.NewDBGatewayEvent(pendingEvent)
	_, err := s.db.NamedExecContext(ctx, "INSERT INTO gateway_pending_events (contract_address, block_number, transaction_index, log_index, block, transaction, type, id, version, new_owner, old_owner, new_location, old_location, new_altitude, old_altitude, new_frequency_plan, old_frequency_plan, new_antenna_gain, old_antenna_gain, time) VALUES (:contract_address, :block_number, :transaction_index, :log_index, :block, :transaction, :type, :id, :version, :new_owner, :old_owner, :new_location, :old_location, :new_altitude, :old_altitude, :new_frequency_plan, :old_frequency_plan, :new_antenna_gain, :old_antenna_gain, :time)", dbPendingGatewayEvent)
	return err
}
