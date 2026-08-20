package controllers

import (
	"errors"
	"net/http"
	"sis-api/config"
	"sis-api/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func parseDateTime(dateStr string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", dateStr)
}

// CREATE
func CreateBorrowing(c *gin.Context) {
	var input models.BorrowingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.Item
	if err := config.DB.WithContext(c.Request.Context()).First(&item, input.ItemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengecek data item"})
		return
	}

	borrowDate, err := parseDateTime(input.BorrowDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format borrow_date tidak valid (Gunakan YYYY-MM-DD HH:mm:ss atau YYYY-MM-DD)"})
		return
	}

	dueDate, err := parseDateTime(input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format due_date tidak valid (Gunakan YYYY-MM-DD HH:mm:ss atau YYYY-MM-DD)"})
		return
	}

	var returnDate *time.Time
	if input.ReturnDate != nil && *input.ReturnDate != "" {
		parsedReturn, err := parseDateTime(*input.ReturnDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format return_date tidak valid"})
			return
		}
		returnDate = &parsedReturn
	}

	conditionBefore := input.ConditionBefore
	if conditionBefore == "" {
		conditionBefore = item.Condition
	}

	if input.Status == "" {
		input.Status = models.StatusPending
	}

	borrowing := models.Borrowing{
		ItemID:          input.ItemID,
		UserID:          input.UserID,
		BorrowDate:      borrowDate,
		DueDate:         dueDate,
		ReturnDate:      returnDate,
		ConditionBefore: conditionBefore,
		ConditionAfter:  input.ConditionAfter,
		Status:          input.Status,
		Notes:           input.Notes,
	}

	if err := config.DB.WithContext(c.Request.Context()).Create(&borrowing).Error; err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan peminjaman"})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("Item").Preload("User").First(&borrowing, borrowing.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Berhasil mengajukan peminjaman",
		"borrowing": borrowing,
	})
}

// READ ALL
func GetAllBorrowings(c *gin.Context) {
	var borrowings []models.Borrowing

	if err := config.DB.WithContext(c.Request.Context()).Preload("Item").Preload("User").Find(&borrowings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"borrowings": borrowings,
	})
}

// READ DETAIL (BY ID)
func GetBorrowingByID(c *gin.Context) {
	id := c.Param("id")
	var borrowing models.Borrowing

	if err := config.DB.WithContext(c.Request.Context()).Preload("Item").Preload("User").First(&borrowing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data peminjaman tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"borrowing": borrowing,
	})
}

// UPDATE
func UpdateBorrowing(c *gin.Context) {
	id := c.Param("id")
	var borrowing models.Borrowing

	if err := config.DB.WithContext(c.Request.Context()).First(&borrowing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data peminjaman tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari data peminjaman"})
		return
	}

	var input models.BorrowingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrowDate, err := parseDateTime(input.BorrowDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format borrow_date tidak valid"})
		return
	}

	dueDate, err := parseDateTime(input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format due_date tidak valid"})
		return
	}

	var returnDate *time.Time
	if input.ReturnDate != nil && *input.ReturnDate != "" {
		parsedReturn, err := parseDateTime(*input.ReturnDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format return_date tidak valid"})
			return
		}
		returnDate = &parsedReturn
	}

	borrowing.ItemID = input.ItemID
	borrowing.UserID = input.UserID
	borrowing.BorrowDate = borrowDate
	borrowing.DueDate = dueDate
	borrowing.ReturnDate = returnDate
	borrowing.ConditionBefore = input.ConditionBefore
	borrowing.ConditionAfter = input.ConditionAfter
	borrowing.Status = input.Status
	borrowing.Notes = input.Notes

	if err := config.DB.WithContext(c.Request.Context()).Save(&borrowing).Error; err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID atau User ID tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data peminjaman"})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("Item").Preload("User").First(&borrowing, borrowing.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Berhasil memperbarui data peminjaman",
		"borrowing": borrowing,
	})
}

// DELETE
func DeleteBorrowing(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.WithContext(c.Request.Context()).Delete(&models.Borrowing{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data peminjaman"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data peminjaman tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus data peminjaman",
	})
}

// APPROVAL PEMINJAMAN (Pending -> Disetujui / Ditolak / Dipinjam)
func ApproveBorrowing(c *gin.Context) {
	id := c.Param("id")

	var input models.BorrowingApprovalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != models.StatusDisetujui && input.Status != models.StatusDitolak && input.Status != models.BorrowStatusDipinjam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status approval harus: disetujui, ditolak, atau dipinjam"})
		return
	}

	var borrowing models.Borrowing
	if err := config.DB.WithContext(c.Request.Context()).Preload("Item").First(&borrowing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data peminjaman tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari data peminjaman"})
		return
	}

	if borrowing.Status != models.StatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Peminjaman ini sudah diproses sebelumnya"})
		return
	}

	err := config.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if input.Status == models.StatusDisetujui || input.Status == models.BorrowStatusDipinjam {
			if borrowing.Item.Status != models.StatusTersedia {
				return errors.New("barang sedang tidak tersedia untuk dipinjam")
			}

			// Menggunakan StatusDipinjam langsung dari enum ItemStatus
			if err := tx.Model(&models.Item{}).Where("id = ?", borrowing.ItemID).Update("status", models.StatusDipinjam).Error; err != nil {
				return err
			}
		}

		if input.Notes != nil {
			borrowing.Notes = input.Notes
		}
		borrowing.Status = input.Status

		return tx.Save(&borrowing).Error
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("Item").Preload("User").First(&borrowing, borrowing.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Status peminjaman berhasil diperbarui",
		"borrowing": borrowing,
	})
}

// PENERIMAAN PENGEMBALIAN BARANG (Dipinjam -> Dikembalikan)
func ReturnBorrowing(c *gin.Context) {
	id := c.Param("id")

	var input models.BorrowingReturnInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var borrowing models.Borrowing
	if err := config.DB.WithContext(c.Request.Context()).Preload("Item").First(&borrowing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data peminjaman tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari data peminjaman"})
		return
	}

	if borrowing.Status != models.BorrowStatusDipinjam && borrowing.Status != models.StatusDisetujui && borrowing.Status != models.StatusTerlambat {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Peminjaman ini tidak dalam status dipinjam"})
		return
	}

	now := time.Now()

	err := config.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		newItemStatus := models.StatusTersedia
		if input.ConditionAfter == models.ConditionRusakBerat {
			newItemStatus = models.StatusPerbaikan
		}

		if err := tx.Model(&models.Item{}).Where("id = ?", borrowing.ItemID).Updates(map[string]interface{}{
			"status":    newItemStatus,
			"condition": input.ConditionAfter,
		}).Error; err != nil {
			return err
		}

		borrowing.ReturnDate = &now
		borrowing.ConditionAfter = &input.ConditionAfter
		borrowing.Status = models.StatusDikembalikan
		if input.Notes != nil {
			borrowing.Notes = input.Notes
		}

		return tx.Save(&borrowing).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses pengembalian barang"})
		return
	}

	config.DB.WithContext(c.Request.Context()).Preload("Item").Preload("User").First(&borrowing, borrowing.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Barang berhasil dikembalikan",
		"borrowing": borrowing,
	})
}
