package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func hallsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listHalls(w, r)
	case http.MethodPost:
		createHall(w, r)
	case http.MethodPut:
		updateHall(w, r)
	case http.MethodDelete:
		deleteHall(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func listHalls(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, theme, max_capacity, status FROM halls ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	halls := []Hall{}
	for rows.Next() {
		var h Hall
		if err := rows.Scan(&h.ID, &h.Name, &h.Theme, &h.MaxCapacity, &h.Status); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		halls = append(halls, h)
	}
	writeJSON(w, http.StatusOK, halls)
}

func createHall(w http.ResponseWriter, r *http.Request) {
	var h Hall
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if h.Name == "" || h.Theme == "" || h.MaxCapacity <= 0 {
		writeError(w, http.StatusBadRequest, "invalid parameters")
		return
	}
	validThemes := map[string]bool{"历史人物": true, "影视明星": true, "体育冠军": true, "世界领袖": true}
	if !validThemes[h.Theme] {
		writeError(w, http.StatusBadRequest, "invalid theme")
		return
	}
	if h.Status == "" {
		h.Status = "开放"
	}
	validStatus := map[string]bool{"开放": true, "关闭": true, "维修": true}
	if !validStatus[h.Status] {
		writeError(w, http.StatusBadRequest, "invalid status")
		return
	}

	res, err := db.Exec("INSERT INTO halls (name, theme, max_capacity, status) VALUES (?, ?, ?, ?)",
		h.Name, h.Theme, h.MaxCapacity, h.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.ID, _ = res.LastInsertId()
	writeJSON(w, http.StatusCreated, h)
}

func updateHall(w http.ResponseWriter, r *http.Request) {
	var h Hall
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if h.ID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	_, err := db.Exec("UPDATE halls SET name=?, theme=?, max_capacity=?, status=? WHERE id=?",
		h.Name, h.Theme, h.MaxCapacity, h.Status, h.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, h)
}

func deleteHall(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	_, err = db.Exec("DELETE FROM halls WHERE id=?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func guidesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listGuides(w, r)
	case http.MethodPost:
		createGuide(w, r)
	case http.MethodPut:
		updateGuide(w, r)
	case http.MethodDelete:
		deleteGuide(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func listGuides(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nickname, phone, language, status FROM guides ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	guides := []Guide{}
	for rows.Next() {
		var g Guide
		if err := rows.Scan(&g.ID, &g.Nickname, &g.Phone, &g.Language, &g.Status); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		guides = append(guides, g)
	}
	writeJSON(w, http.StatusOK, guides)
}

func createGuide(w http.ResponseWriter, r *http.Request) {
	var g Guide
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if g.Nickname == "" || g.Phone == "" || g.Language == "" {
		writeError(w, http.StatusBadRequest, "invalid parameters")
		return
	}
	validLangs := map[string]bool{"普通话": true, "粤语": true, "英语": true, "日语": true}
	if !validLangs[g.Language] {
		writeError(w, http.StatusBadRequest, "invalid language")
		return
	}
	if g.Status == "" {
		g.Status = "在职"
	}

	res, err := db.Exec("INSERT INTO guides (nickname, phone, language, status) VALUES (?, ?, ?, ?)",
		g.Nickname, g.Phone, g.Language, g.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	g.ID, _ = res.LastInsertId()
	writeJSON(w, http.StatusCreated, g)
}

func updateGuide(w http.ResponseWriter, r *http.Request) {
	var g Guide
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if g.ID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	_, err := db.Exec("UPDATE guides SET nickname=?, phone=?, language=?, status=? WHERE id=?",
		g.Nickname, g.Phone, g.Language, g.Status, g.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func deleteGuide(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	_, err = db.Exec("DELETE FROM guides WHERE id=?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func reservationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listReservations(w, r)
	case http.MethodPost:
		createReservation(w, r)
	case http.MethodDelete:
		deleteReservation(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func listReservations(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT r.id, r.guest_name, r.guest_phone, r.hall_id, h.name, r.time_slot, 
		       r.people_count, r.guide_id, g.nickname, r.created_at
		FROM reservations r
		LEFT JOIN halls h ON r.hall_id = h.id
		LEFT JOIN guides g ON r.guide_id = g.id
		ORDER BY r.time_slot DESC, r.id DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	list := []Reservation{}
	for rows.Next() {
		var rv Reservation
		var guideID sql.NullInt64
		var guideName sql.NullString
		var hallName sql.NullString
		if err := rows.Scan(&rv.ID, &rv.GuestName, &rv.GuestPhone, &rv.HallID, &hallName,
			&rv.TimeSlot, &rv.PeopleCount, &guideID, &guideName, &rv.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if hallName.Valid {
			rv.HallName = hallName.String
		}
		if guideID.Valid {
			id := guideID.Int64
			rv.GuideID = &id
		}
		if guideName.Valid {
			rv.GuideName = guideName.String
		}
		list = append(list, rv)
	}
	writeJSON(w, http.StatusOK, list)
}

func createReservation(w http.ResponseWriter, r *http.Request) {
	var rv Reservation
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if rv.GuestName == "" || rv.GuestPhone == "" || rv.HallID <= 0 || rv.TimeSlot == "" || rv.PeopleCount <= 0 {
		writeError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	var status string
	var maxCap int
	err = tx.QueryRow("SELECT status, max_capacity FROM halls WHERE id=?", rv.HallID).Scan(&status, &maxCap)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusBadRequest, "hall not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if status != "开放" {
		writeError(w, http.StatusBadRequest, "hall is not open for reservation")
		return
	}

	var currentCount int
	err = tx.QueryRow("SELECT COALESCE(SUM(people_count), 0) FROM reservations WHERE hall_id=? AND time_slot=?",
		rv.HallID, rv.TimeSlot).Scan(&currentCount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if currentCount+rv.PeopleCount > maxCap {
		writeError(w, http.StatusBadRequest, "exceeds hall max capacity")
		return
	}

	res, err := tx.Exec(`INSERT INTO reservations (guest_name, guest_phone, hall_id, time_slot, people_count) 
		VALUES (?, ?, ?, ?, ?)`, rv.GuestName, rv.GuestPhone, rv.HallID, rv.TimeSlot, rv.PeopleCount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	rv.ID, _ = res.LastInsertId()

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, rv)
}

func deleteReservation(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	_, err = db.Exec("DELETE FROM reservations WHERE id=?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func assignHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		ReservationID int64 `json:"reservation_id"`
		GuideID       int64 `json:"guide_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ReservationID <= 0 || req.GuideID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	var guideStatus string
	err = tx.QueryRow("SELECT status FROM guides WHERE id=?", req.GuideID).Scan(&guideStatus)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusBadRequest, "guide not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if guideStatus != "在职" {
		writeError(w, http.StatusBadRequest, "guide is not active")
		return
	}

	var timeSlot string
	err = tx.QueryRow("SELECT time_slot FROM reservations WHERE id=?", req.ReservationID).Scan(&timeSlot)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusBadRequest, "reservation not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var conflictCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM reservations 
		WHERE guide_id=? AND time_slot=? AND id<>?`, req.GuideID, timeSlot, req.ReservationID).Scan(&conflictCount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if conflictCount > 0 {
		writeError(w, http.StatusBadRequest, "guide already assigned at this time slot")
		return
	}

	_, err = tx.Exec("UPDATE reservations SET guide_id=? WHERE id=?", req.GuideID, req.ReservationID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "assigned"})
}

func monthlyThemeStatsHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	loc := now.Location()
	year, month, _ := now.Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	startOfNext := startOfMonth.AddDate(0, 1, 0)

	rows, err := db.Query(`
		SELECT h.theme, COUNT(r.id) 
		FROM halls h 
		LEFT JOIN reservations r ON h.id = r.hall_id 
			AND r.created_at >= ? AND r.created_at < ?
		GROUP BY h.theme
	`, startOfMonth.Format("2006-01-02 15:04:05"), startOfNext.Format("2006-01-02 15:04:05"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	stats := []ThemeStat{}
	for rows.Next() {
		var s ThemeStat
		if err := rows.Scan(&s.Theme, &s.Count); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		stats = append(stats, s)
	}
	writeJSON(w, http.StatusOK, stats)
}

func todayReservationsHandler(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")
	prefix := today + "%"
	query := `
		SELECT r.id, r.guest_name, r.guest_phone, r.hall_id, h.name, r.time_slot, 
		       r.people_count, r.guide_id, g.nickname, r.created_at
		FROM reservations r
		LEFT JOIN halls h ON r.hall_id = h.id
		LEFT JOIN guides g ON r.guide_id = g.id
		WHERE r.time_slot LIKE ?
		ORDER BY r.time_slot ASC, r.id ASC
	`
	rows, err := db.Query(query, prefix)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	list := []Reservation{}
	for rows.Next() {
		var rv Reservation
		var guideID sql.NullInt64
		var guideName sql.NullString
		var hallName sql.NullString
		if err := rows.Scan(&rv.ID, &rv.GuestName, &rv.GuestPhone, &rv.HallID, &hallName,
			&rv.TimeSlot, &rv.PeopleCount, &guideID, &guideName, &rv.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if hallName.Valid {
			rv.HallName = hallName.String
		}
		if guideID.Valid {
			id := guideID.Int64
			rv.GuideID = &id
		}
		if guideName.Valid {
			rv.GuideName = guideName.String
		}
		list = append(list, rv)
	}
	_ = strings.TrimSpace
	writeJSON(w, http.StatusOK, list)
}
