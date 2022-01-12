package ast

func (d *Document) Clone() *Document {
	input := d.Input.Clone()

	rootNodes := make([]Node, len(d.RootNodes))
	copy(rootNodes, d.RootNodes)

	schemaDefinitions := make([]SchemaDefinition, len(d.SchemaDefinitions))
	copy(schemaDefinitions, d.SchemaDefinitions)

	schemaExtensions := make([]SchemaExtension, len(d.SchemaExtensions))


	rootOperationTypeDefinitions := make([]RootOperationTypeDefinition, len(d.RootOperationTypeDefinitions))
	copy(rootOperationTypeDefinitions, d.RootOperationTypeDefinitions)



}
