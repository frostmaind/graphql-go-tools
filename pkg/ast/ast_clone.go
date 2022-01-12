package ast

func (d *Document) Clone() *Document {
	input := d.Input.Clone()

	rootNodes := make([]Node, len(d.RootNodes))
	copy(rootNodes, d.RootNodes)

	schemaDefinitions := make([]SchemaDefinition, len(d.SchemaDefinitions))
	copy(schemaDefinitions, d.SchemaDefinitions)

	schemaExtensions := make([]SchemaExtension, len(d.SchemaExtensions))
	for i := range d.SchemaExtensions {
		schemaExtensions[i] = d.SchemaExtensions[i].Clone()
	}

	rootOperationTypeDefinitions := make([]RootOperationTypeDefinition, len(d.RootOperationTypeDefinitions))
	copy(rootOperationTypeDefinitions, d.RootOperationTypeDefinitions)

	directives := make([]Directive, len(d.Directives))
	for i := range d.Directives {
		directives[i] = d.Directives[i].Clone()
	}

	arguments := make([]Argument, len(d.Arguments))
	copy(arguments, d.Arguments)

	objectTypeDefinitions := make([]ObjectTypeDefinition, len(d.ObjectTypeDefinitions))
	for i := range d.ObjectTypeDefinitions {
		objectTypeDefinitions[i] = d.ObjectTypeDefinitions[i].Clone()
	}

	objectTypeExtensions := make([]ObjectTypeExtension, len(d.ObjectTypeExtensions))
	for i := range d.ObjectTypeExtensions {
		objectTypeExtensions[i] = d.ObjectTypeExtensions[i].Clone()
	}

	fieldDefinitions := make([]FieldDefinition, len(d.FieldDefinitions))
	for i := range d.FieldDefinitions {
		fieldDefinitions[i] = d.FieldDefinitions[i].Clone()
	}

	types := make([]Type, len(d.Types))
	copy(types, d.Types)

	inputValueDefinitions := make([]InputValueDefinition, len(d.InputValueDefinitions))
	for i := range d.InputValueDefinitions {
		inputValueDefinitions[i] = d.InputValueDefinitions[i].Clone()
	}

	inputObjectTypeDefinitions := make([]InputObjectTypeDefinition, len(d.InputObjectTypeDefinitions))
	for i := range d.InputObjectTypeDefinitions {
		inputObjectTypeDefinitions[i] = d.InputObjectTypeDefinitions[i].Clone()
	}

	inputObjectTypeExtensions := make([]InputObjectTypeExtension, len(d.InputObjectTypeExtensions))
	for i := range d.InputObjectTypeExtensions {
		inputObjectTypeExtensions[i] = d.InputObjectTypeExtensions[i].Clone()
	}

	scalarTypeDefinitions := make([]ScalarTypeDefinition, len(d.ScalarTypeDefinitions))
	for i := range d.ScalarTypeDefinitions {
		scalarTypeDefinitions[i] = d.ScalarTypeDefinitions[i].Clone()
	}

	scalarTypeExtensions := make([]ScalarTypeExtension, len(d.ScalarTypeExtensions))
	for i := range d.ScalarTypeExtensions {
		scalarTypeExtensions[i] = d.ScalarTypeExtensions[i].Clone()
	}

	interfaceTypeDefinitions := make([]InterfaceTypeDefinition, len(d.InterfaceTypeDefinitions))
	for i := range d.InterfaceTypeDefinitions {
		interfaceTypeDefinitions[i] = d.InterfaceTypeDefinitions[i].Clone()
	}

	interfaceTypeExtensions := make([]InterfaceTypeExtension, len(d.InterfaceTypeExtensions))
	for i := range d.InterfaceTypeExtensions {
		interfaceTypeExtensions[i] = d.InterfaceTypeExtensions[i].Clone()
	}

	unionTypeDefinitions := make([]UnionTypeDefinition, len(d.UnionTypeDefinitions))
	for i := range d.UnionTypeDefinitions {
		unionTypeDefinitions[i] = d.UnionTypeDefinitions[i].Clone()
	}

	unionTypeExtensions := make([]UnionTypeExtension, len(d.UnionTypeExtensions))
	for i := range d.UnionTypeExtensions {
		unionTypeExtensions[i] = d.UnionTypeExtensions[i].Clone()
	}

	enumTypeDefinitions := make([]EnumTypeDefinition, len(d.EnumTypeDefinitions))
	for i := range d.EnumTypeDefinitions {
		enumTypeDefinitions[i] = d.EnumTypeDefinitions[i].Clone()
	}

	enumTypeExtensions := make([]EnumTypeExtension, len(d.EnumTypeExtensions))
	for i := range d.EnumTypeExtensions {
		enumTypeExtensions[i] = d.EnumTypeExtensions[i].Clone()
	}

	enumValueDefinitions := make([]EnumValueDefinition, len(d.EnumValueDefinitions))
	for i := range d.EnumValueDefinitions {
		enumValueDefinitions[i] = d.EnumValueDefinitions[i].Clone()
	}

	directiveDefinitions := make([]DirectiveDefinition, len(d.DirectiveDefinitions))
	for i := range d.DirectiveDefinitions {
		directiveDefinitions[i] = d.DirectiveDefinitions[i].Clone()
	}

	values := make([]Value, len(d.Values))
	copy(values, d.Values)

	listValues := make([]ListValue, len(d.ListValues))
	for i := range d.ListValues {
		listValues[i] = d.ListValues[i].Clone()
	}

	variableValues := make([]VariableValue, len(d.VariableValues))
	copy(variableValues, d.VariableValues)

	stringValues := make([]StringValue, len(d.StringValues))
	copy(stringValues, d.StringValues)

	intValues := make([]IntValue, len(d.IntValues))
	copy(intValues, d.IntValues)

	floatValues := make([]FloatValue, len(d.FloatValues))
	copy(floatValues, d.FloatValues)

	enumValues := make([]EnumValue, len(d.EnumValues))
	copy(enumValues, d.EnumValues)

	objectFields := make([]ObjectField, len(d.ObjectFields))
	copy(objectFields, d.ObjectFields)

	objectValues := make([]ObjectValue, len(d.ObjectValues))
	for i := range d.ObjectValues {
		objectValues[i] = d.ObjectValues[i].Clone()
	}

	selections := make([]Selection, len(d.Selections))
	copy(selections, d.Selections)

	selectionSets := make([]SelectionSet, len(d.SelectionSets))
	for i := range d.SelectionSets {
		selectionSets[i] = d.SelectionSets[i].Clone()
	}

	fields := make([]Field, len(d.Fields))
	for i := range d.Fields {
		fields[i] = d.Fields[i].Clone()
	}

	inlineFragments := make([]InlineFragment, len(d.InlineFragments))
	for i := range d.InlineFragments {
		inlineFragments[i] = d.InlineFragments[i].Clone()
	}

	fragmentSpreads := make([]FragmentSpread, len(d.FragmentSpreads))
	for i := range d.FragmentSpreads {
		fragmentSpreads[i] = d.FragmentSpreads[i].Clone()
	}

	operationDefinitions := make([]OperationDefinition, len(d.OperationDefinitions))
	for i := range d.OperationDefinitions {
		operationDefinitions[i] = d.OperationDefinitions[i].Clone()
	}

	variableDefinitions := make([]VariableDefinition, len(d.VariableDefinitions))
	for i := range d.VariableDefinitions {
		variableDefinitions[i] = d.VariableDefinitions[i].Clone()
	}

	fragmentDefinitions := make([]FragmentDefinition, len(d.FragmentDefinitions))
	for i := range d.FragmentDefinitions {
		fragmentDefinitions[i] = d.FragmentDefinitions[i].Clone()
	}

	booleanValues := d.BooleanValues

	refs := make([][8]int, len(d.Refs))
	for i := range d.Refs {
		refs[i] = d.Refs[i]
	}

	refIndex := d.RefIndex
	index := d.Index.Clone()

	return &Document{
		Input:                        input,
		RootNodes:                    rootNodes,
		SchemaDefinitions:            schemaDefinitions,
		SchemaExtensions:             schemaExtensions,
		RootOperationTypeDefinitions: rootOperationTypeDefinitions,
		Directives:                   directives,
		Arguments:                    arguments,
		ObjectTypeDefinitions:        objectTypeDefinitions,
		ObjectTypeExtensions:         objectTypeExtensions,
		FieldDefinitions:             fieldDefinitions,
		Types:                        types,
		InputValueDefinitions:        inputValueDefinitions,
		InputObjectTypeDefinitions:   inputObjectTypeDefinitions,
		InputObjectTypeExtensions:    inputObjectTypeExtensions,
		ScalarTypeDefinitions:        scalarTypeDefinitions,
		ScalarTypeExtensions:         scalarTypeExtensions,
		InterfaceTypeDefinitions:     interfaceTypeDefinitions,
		InterfaceTypeExtensions:      interfaceTypeExtensions,
		UnionTypeDefinitions:         unionTypeDefinitions,
		UnionTypeExtensions:          unionTypeExtensions,
		EnumTypeDefinitions:          enumTypeDefinitions,
		EnumTypeExtensions:           enumTypeExtensions,
		EnumValueDefinitions:         enumValueDefinitions,
		DirectiveDefinitions:         directiveDefinitions,
		Values:                       values,
		ListValues:                   listValues,
		VariableValues:               variableValues,
		StringValues:                 stringValues,
		IntValues:                    intValues,
		FloatValues:                  floatValues,
		EnumValues:                   enumValues,
		ObjectFields:                 objectFields,
		ObjectValues:                 objectValues,
		Selections:                   selections,
		SelectionSets:                selectionSets,
		Fields:                       fields,
		InlineFragments:              inlineFragments,
		FragmentSpreads:              fragmentSpreads,
		OperationDefinitions:         operationDefinitions,
		VariableDefinitions:          variableDefinitions,
		FragmentDefinitions:          fragmentDefinitions,
		BooleanValues:                booleanValues,
		Refs:                         refs,
		RefIndex:                     refIndex,
		Index:                        index,
	}
}
