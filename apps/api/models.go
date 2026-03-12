package main 

import "time" 

// User representa a tabela 'users'
type User struct {
		ID			 string		`json:"id"`
		Alias		 string		`json:"alias"`
		Email		 string		`json:"email"`	
		PasswordHash string		`json:"-"`	
		CurrentLevel int		`json:"current_level"`// O traço significa: NUNCA envie a senha no JSON de resposta
		CurrentXP	 int		`json:"current_xp"`	
		Credits		 int		`json:"credits"`	
		Designation	 string		`json:"designation"`
		CreatedAt    time.Time `json:"created_at"`			
}

// Hobby representa a tabela 'hobbies'
type Hobby struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	CurrentXP   int    `json:"current_xp"`
	NextLevelXP int    `json:"next_level_xp"`
	IconName    string `json:"icon_name"`
}

// Protocol representa os hábitos diários (tabela 'protocols')
type Protocol struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Title           string    `json:"title"`
	StreakCount     int       `json:"streak_count"`
	AttributeType   string    `json:"attribute_type"`
	LastCompletedAt time.Time `json:"last_completed_at"`
	IsActive        bool      `json:"is_active"`
}

// Operation representa as tarefas únicas (tabela 'operations')
type Operation struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	HobbyID      *string   `json:"hobby_id"` // Ponteiro porque pode ser nulo (tarefa sem hobby)
	Title        string    `json:"title"`
	Priority     string    `json:"priority"`
	XPReward     int       `json:"xp_reward"`
	CreditReward int       `json:"credit_reward"`
	IsCompleted  bool      `json:"is_completed"`
	Deadline     time.Time `json:"deadline"`
	CompletedAt  time.Time `json:"completed_at"`
}

// LedgerEntry representa o extrato de XP (tabela 'ledger_entries')
type LedgerEntry struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      int       `json:"amount"`
	EntryType   string    `json:"entry_type"` // GAINED, DECAY, SPENT
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}