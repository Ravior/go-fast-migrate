package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/Ravior/go-fast-migrate/migrations"
	"github.com/Ravior/go-fast-migrate/util"
	"github.com/golang-module/carbon"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	createMigrationSql = "CREATE TABLE IF NOT EXISTS migrations (" +
		"id int(10) UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY," +
		"migration varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"batch int(11) NOT NULL " +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"
	updateMigrationSql    = "INSERT INTO migrations (migration, batch) VALUES DummyString;"
	queryAllMigrationSql  = "SELECT * FROM migrations;"
	queryLastMigrationSql = "SELECT batch FROM migrations ORDER BY batch DESC;"
	migrationPath         = "./migrations/" // MigrationInte files save path
	migrationTemplate     = "./template/migrate.template"
	conn                  = util.ConfigHelper.GetString("db.migration.db", "default") // 迁移对应的数据库链接Key
)

func init() {
	initMigration()
}

func initMigration() {
	// CreateMigration migrations table if not exist
	_, err := util.DbHelper.Exec(createMigrationSql, conn)
	util.SysHelper.CheckErr(err)
}

func checkArgs() {
	if len(os.Args) < 2 {
		util.SysHelper.Exit("Please Set Comamnd")
	}
}

func main() {
	util.DbHelper.GetConn("default")
	checkArgs()

	command := os.Args[1]

	if strings.Compare(command, "new") == 0 {
		fileName := ""
		if len(os.Args) > 2 {
			fileName = os.Args[2]
		}

		if len(fileName) <= 0 {
			util.SysHelper.Exit("Please enter a migration file name")
		}

		util.LogHelper.Info("Try to create migration file: %s", fileName)

		path, err := CreateMigration(fileName)

		if err != nil {
			util.SysHelper.Exit("创建迁移文件失败，错误信息：%v", err)
		}

		util.LogHelper.Info("CreateTable migration success! Path: %s", path)

	} else if strings.Compare(command, "start") == 0 {
		err := Migrate()
		if err != nil {
			util.SysHelper.Exit("执行迁移任务失败，错误信息：%v", err)
		}
		util.LogHelper.Info("执行迁移任务成功")
	} else if strings.Compare(command, "refresh") == 0 {
		_, err := Refresh()
		if err != nil {
			util.SysHelper.Exit("执行刷新任务失败，错误信息：%v", err)
		}
	} else if strings.Compare(command, "rollback") == 0 {
		var step string
		if len(os.Args) < 3 {
			step = "1"
		} else {
			step = os.Args[2]
		}
		err := Rollback(step)
		if err != nil {
			util.SysHelper.Exit("执行迁移任务失败，错误信息：%v", err)
		}
	} else {
		util.SysHelper.Exit("Command not support: %v", command)
	}

}

// 创建新的迁移文件
func CreateMigration(name string) (string, error) {

	timestamp := carbon.Now().Format("Y_m_d_His")

	filename := fmt.Sprintf("%s_%s", timestamp, name)
	filePath := fmt.Sprintf("%s%s.go", migrationPath, filename)

	migrateFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer migrateFile.Close()

	template, err := util.FileHelper.ReadFile(migrationTemplate)

	if err != nil {
		return "", err
	}

	// 迁移模板文件
	templateStr := string(template)

	nameArr := strings.Split(name, "_")
	newNameArr := make([]string, len(nameArr))
	for _, key := range nameArr {
		newNameArr = append(newNameArr, util.StrHelper.Ucfirst(key))
	}

	newName := strings.Join(newNameArr, "") + timestamp

	migrateFileWriter := bufio.NewWriter(migrateFile)

	// 替换模板中占位字符串
	rs := strings.Replace(templateStr, "DummyString", newName, -1)
	rs2 := strings.Replace(rs, "dummyString", util.StrHelper.Lcfirst(newName), -1)
	rs3 := strings.Replace(rs2, "fileString", filename, -1)

	_, err = migrateFileWriter.WriteString(rs3)

	if err != nil {
		return "", err
	}

	_ = migrateFileWriter.Flush()

	return filePath, nil
}

// 执行迁移
func Migrate() error {
	var (
		fSlices          []string // 迁移目录下的文件列表
		batch            int
		lastBatch        int
		hasMigratedFiles []string // 已经执行过迁移的文件
		toMigrateFiles   []string // 等待执行迁移的文件
		insertStr        string
		symbol           string
	)

	// List migrations files
	files, err := ioutil.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		filename := strings.Replace(f.Name(), ".go", "", -1)
		match, err := regexp.MatchString("\\d{4}_\\d{2}_\\d{2}_(\\w)+", filename)
		if err != nil {
			util.LogHelper.Warn("UnKown Migrate File: %s, Error: %v", filename, err)
			continue
		}
		if match {
			fSlices = append(fSlices, filename)
		}

	}

	// Check migration version in database
	rows, err := util.DbHelper.Query(queryAllMigrationSql, conn)
	if err != nil {
		return err
	}

	lastRow := util.DbHelper.QueryRow(queryLastMigrationSql, conn)
	_ = lastRow.Scan(&lastBatch)
	// 计算本次执行的批次号
	batch = lastBatch + 1

	defer rows.Close()

	if lastBatch == 0 {
		// No migration record in database, all migrations should to be Migrate
		toMigrateFiles = fSlices
	} else {
		// Get migrated files' name
		for rows.Next() {
			// Row to Migration Struct
			m, err := scanRow(rows)
			if err != nil {
				return err
			}

			hasMigratedFiles = append(hasMigratedFiles, m.Migration)
		}

		// Compare and get which migration not migrated yet
		for _, v := range fSlices {
			if !util.ArrHelper.StrArrContain(hasMigratedFiles, v) {
				toMigrateFiles = append(toMigrateFiles, v)
			}
		}
	}

	// Nothing to Migrate, stop and log fatal
	toMigrateLen := len(toMigrateFiles)
	if toMigrateLen == 0 {
		util.SysHelper.Exit("暂无可执行迁移计划")
	}

	// Migrate
	for i, name := range toMigrateFiles {
		migrations.DataMap[name].Up()

		// Calculate the batch number, which is need to Migrate
		if i+1 == toMigrateLen {
			symbol = ""
		} else {
			symbol = ","
		}

		insertStr += "('" + name + "', " + strconv.Itoa(batch) + ")" + symbol
	}

	// Connect sql update statement
	updateMigrationSql = strings.Replace(updateMigrationSql, "DummyString", insertStr, -1)

	_, err = util.DbHelper.Exec(updateMigrationSql, conn)
	if err != nil {
		return err
	}

	return nil
}

func Rollback(step string) error {
	var (
		lastBatch   int
		toBatch     int
		err         error
		rows        *sql.Rows
		m           *Migration
		rollBackMig []string
	)

	lastRow := util.DbHelper.QueryRow(queryLastMigrationSql, conn)
	_ = lastRow.Scan(&lastBatch)

	if i, err := strconv.Atoi(step); err == nil {
		if lastBatch >= i {
			toBatch = lastBatch - (i - 1)
		} else {
			util.LogHelper.Error("Nothing to rollback")
			return err
		}
	}

	// Which migrations need to be Rollback
	rows, err = util.DbHelper.Query("SELECT * FROM migrations WHERE `batch`>="+strconv.Itoa(toBatch), conn)
	if err != nil {
		return err
	}

	for rows.Next() {
		m, err = scanRow(rows)
		if err != nil {
			return err
		}

		rollBackMig = append(rollBackMig, m.Migration)
	}

	for _, name := range rollBackMig {
		migrations.DataMap[name].Down()
	}

	// Delete migrations record
	_, err = util.DbHelper.Exec("DELETE FROM migrations WHERE `batch`>=" + strconv.Itoa(toBatch))
	if err != nil {
		return err
	}

	return nil
}

// Refresh migration: Rollback all and re-Migrate
func Refresh() (bool, error) {
	var (
		insertStr   string
		symbol      string
		err         error
		rows        *sql.Rows
		rollBackMig []string
		m           *Migration
	)

	rows, err = util.DbHelper.Query("SELECT * FROM migrations;", conn)
	if err != nil {
		return false, err
	}

	for rows.Next() {
		m, err = scanRow(rows)
		if err != nil {
			return false, err
		}

		rollBackMig = append(rollBackMig, m.Migration)
	}

	// Rollback and re-Migrate
	fileLen := len(rollBackMig)
	if fileLen > 0 {
		for i, name := range rollBackMig {
			// down
			migrations.DataMap[name].Down()

			// up
			migrations.DataMap[name].Up()

			if i == fileLen-1 {
				symbol = ""
			} else {
				symbol = ","
			}

			insertStr += "('" + name + "', 1)" + symbol
		}

		// Update migrations table
		_, _ = util.DbHelper.Exec("TRUNCATE migrations;", conn)
		_, err = util.DbHelper.Exec(strings.Replace(updateMigrationSql, "DummyString", insertStr, -1), conn)
		if err != nil {
			return false, err
		}

		return true, nil

	} else {
		return false, nil
	}
}

type Migration struct {
	ID        int64
	Migration string
	Batch     int64
}

type rowScanner interface {
	Scan(dst ...interface{}) error
}

// Map sql row to struct
func scanRow(s rowScanner) (*Migration, error) {
	var (
		id        int64
		migration sql.NullString
		batch     int64
	)

	if err := s.Scan(&id, &migration, &batch); err != nil {
		return nil, err
	}

	mig := &Migration{
		ID:        id,
		Migration: migration.String,
		Batch:     batch,
	}

	return mig, nil
}
