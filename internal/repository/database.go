package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgresDB 创建PostgreSQL数据库连接
func NewPostgresDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 连接池配置
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库Ping失败: %w", err)
	}

	log.Println("数据库连接成功")
	return db, nil
}

// RunMigrations 运行数据库迁移
func RunMigrations(db *sql.DB) error {
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_users",
			sql: `CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username VARCHAR(50) UNIQUE NOT NULL,
				email VARCHAR(100) UNIQUE NOT NULL,
				password VARCHAR(255) NOT NULL,
				avatar TEXT DEFAULT '',
				phone VARCHAR(20) DEFAULT '',
				total_dist DECIMAL(10,2) DEFAULT 0,
				total_rides INT DEFAULT 0,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_routes",
			sql: `CREATE TABLE IF NOT EXISTS routes (
				id SERIAL PRIMARY KEY,
				user_id INT REFERENCES users(id) ON DELETE CASCADE,
				name VARCHAR(100) NOT NULL,
				description TEXT DEFAULT '',
				distance DECIMAL(10,2) NOT NULL,
				duration INT NOT NULL,
				elevation INT DEFAULT 0,
				difficulty VARCHAR(20) NOT NULL,
				coordinates TEXT NOT NULL,
				is_favorite BOOLEAN DEFAULT false,
				is_public BOOLEAN DEFAULT false,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_route_segments",
			sql: `CREATE TABLE IF NOT EXISTS route_segments (
				id SERIAL PRIMARY KEY,
				route_id INT REFERENCES routes(id) ON DELETE CASCADE,
				name VARCHAR(100) NOT NULL,
				distance DECIMAL(10,2) NOT NULL,
				elevation INT DEFAULT 0,
				terrain VARCHAR(50) DEFAULT 'paved',
				tips TEXT DEFAULT '',
				order_num INT NOT NULL
			)`,
		},
		{
			name: "create_equipment",
			sql: `CREATE TABLE IF NOT EXISTS equipment (
				id SERIAL PRIMARY KEY,
				user_id INT REFERENCES users(id) ON DELETE CASCADE,
				name VARCHAR(100) NOT NULL,
				icon VARCHAR(50) DEFAULT 'bike',
				detail TEXT DEFAULT '',
				status VARCHAR(20) NOT NULL,
				category VARCHAR(20) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_consumables",
			sql: `CREATE TABLE IF NOT EXISTS consumables (
				id SERIAL PRIMARY KEY,
				user_id INT REFERENCES users(id) ON DELETE CASCADE,
				name VARCHAR(100) NOT NULL,
				percent INT NOT NULL CHECK (percent >= 0 AND percent <= 100),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_rides",
			sql: `CREATE TABLE IF NOT EXISTS rides (
				id SERIAL PRIMARY KEY,
				user_id INT REFERENCES users(id) ON DELETE CASCADE,
				route_id INT REFERENCES routes(id) ON DELETE SET NULL,
				name VARCHAR(100) NOT NULL,
				distance DECIMAL(10,2) DEFAULT 0,
				duration INT DEFAULT 0,
				avg_speed DECIMAL(5,2) DEFAULT 0,
				max_speed DECIMAL(5,2) DEFAULT 0,
				elevation INT DEFAULT 0,
				calories INT DEFAULT 0,
				start_time TIMESTAMP NOT NULL,
				end_time TIMESTAMP,
				is_active BOOLEAN DEFAULT true,
				is_paused BOOLEAN DEFAULT false,
				coordinates TEXT DEFAULT '',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_ride_points",
			sql: `CREATE TABLE IF NOT EXISTS ride_points (
				id SERIAL PRIMARY KEY,
				ride_id INT REFERENCES rides(id) ON DELETE CASCADE,
				latitude DECIMAL(10,8) NOT NULL,
				longitude DECIMAL(11,8) NOT NULL,
				speed DECIMAL(5,2) DEFAULT 0,
				elevation INT DEFAULT 0,
				timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_teams",
			sql: `CREATE TABLE IF NOT EXISTS teams (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				creator_id INT REFERENCES users(id) ON DELETE CASCADE,
				code VARCHAR(10) UNIQUE NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_team_members",
			sql: `CREATE TABLE IF NOT EXISTS team_members (
				id SERIAL PRIMARY KEY,
				team_id INT REFERENCES teams(id) ON DELETE CASCADE,
				user_id INT REFERENCES users(id) ON DELETE CASCADE,
				joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				color VARCHAR(20) NOT NULL,
				UNIQUE(team_id, user_id)
			)`,
		},
		{
			name: "create_team_messages",
			sql: `CREATE TABLE IF NOT EXISTS team_messages (
				id SERIAL PRIMARY KEY,
				team_id INT REFERENCES teams(id) ON DELETE CASCADE,
				user_id INT REFERENCES users(id) ON DELETE CASCADE,
				content TEXT NOT NULL,
				type VARCHAR(20) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "create_team_locations",
			sql: `CREATE TABLE IF NOT EXISTS team_locations (
				id SERIAL PRIMARY KEY,
				team_id INT REFERENCES teams(id) ON DELETE CASCADE,
				user_id INT REFERENCES users(id) ON DELETE CASCADE,
				latitude DECIMAL(10,8) NOT NULL,
				longitude DECIMAL(11,8) NOT NULL,
				speed DECIMAL(5,2) DEFAULT 0,
				distance DECIMAL(10,2) DEFAULT 0,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(team_id, user_id)
			)`,
		},
	}

	for _, m := range migrations {
		if _, err := db.Exec(m.sql); err != nil {
			return fmt.Errorf("迁移 %s 失败: %w", m.name, err)
		}
		log.Printf("迁移 %s 完成", m.name)
	}

	// 创建索引
	indexes := []struct {
		name string
		sql  string
	}{
		{"idx_routes_user_id", "CREATE INDEX IF NOT EXISTS idx_routes_user_id ON routes(user_id)"},
		{"idx_routes_is_public", "CREATE INDEX IF NOT EXISTS idx_routes_is_public ON routes(is_public)"},
		{"idx_equipment_user_id", "CREATE INDEX IF NOT EXISTS idx_equipment_user_id ON equipment(user_id)"},
		{"idx_equipment_category", "CREATE INDEX IF NOT EXISTS idx_equipment_category ON equipment(category)"},
		{"idx_rides_user_id", "CREATE INDEX IF NOT EXISTS idx_rides_user_id ON rides(user_id)"},
		{"idx_rides_is_active", "CREATE INDEX IF NOT EXISTS idx_rides_is_active ON rides(is_active)"},
		{"idx_rides_start_time", "CREATE INDEX IF NOT EXISTS idx_rides_start_time ON rides(start_time)"},
		{"idx_ride_points_ride_id", "CREATE INDEX IF NOT EXISTS idx_ride_points_ride_id ON ride_points(ride_id)"},
		{"idx_ride_points_timestamp", "CREATE INDEX IF NOT EXISTS idx_ride_points_timestamp ON ride_points(timestamp)"},
		{"idx_team_members_team_id", "CREATE INDEX IF NOT EXISTS idx_team_members_team_id ON team_members(team_id)"},
		{"idx_team_messages_team_id", "CREATE INDEX IF NOT EXISTS idx_team_messages_team_id ON team_messages(team_id)"},
		{"idx_team_locations_team_id", "CREATE INDEX IF NOT EXISTS idx_team_locations_team_id ON team_locations(team_id)"},
	}

	for _, idx := range indexes {
		if _, err := db.Exec(idx.sql); err != nil {
			return fmt.Errorf("创建索引 %s 失败: %w", idx.name, err)
		}
	}

	log.Println("数据库迁移全部完成")
	return nil
}

// InitDB 初始化数据库连接并运行迁移
func InitDB(cfg *Config) (*sql.DB, error) {
	db, err := NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// 应用连接池配置
	db.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)
	db.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DatabaseMaxLifetime)

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// Config 数据库配置（简化版，用于 InitDB）
type Config struct {
	DatabaseURL         string
	DatabaseMaxOpenConns int
	DatabaseMaxIdleConns int
	DatabaseMaxLifetime  time.Duration
}
