package schema

type Schema struct {
	conn string
}

// CreateTable creat a new table
func (s *Schema) CreateTable(table string, apply func(builder *Builder)) {
	builder := NewBuilder(table, s.conn)
	builder.Create()
	apply(builder)

	builder.Run()
}

// AlterTable update a existing table
func (s *Schema) AlterTable(table string, apply func(builder *Builder)) {
	builder := NewBuilder(table, s.conn)
	apply(builder)

	builder.Run()
}

// DropTable drop a table
func (s *Schema) DropTable(table string) {
	builder := NewBuilder(table, s.conn)
	builder.DropTable()

	builder.Run()
}

// DropTable drop a table if exists
func (s *Schema) DropTableIfExists(table string) {
	builder := NewBuilder(table, s.conn)
	builder.DropTableIfExists()

	builder.Run()
}

// NewSchema create a new Schema
func NewSchema(conns ...string) *Schema {
	conn := "default"
	if len(conns) > 0 {
		conn = conns[0]
	}
	return &Schema{
		conn: conn,
	}
}
