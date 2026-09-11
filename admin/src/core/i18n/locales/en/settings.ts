export default {
  settings: {
    title: 'Settings', description: 'Manage platform configuration by namespace with contract-validated value types.',
    unavailable: 'Settings are temporarily unavailable', createTitle: 'Add setting', createDescription: 'Use namespace + key to identify a setting without changing the database schema.',
    namespace: 'Namespace', key: 'Key', type: 'Value type', chooseType: 'Choose a value type', value: 'Value', valuePlaceholder: 'Enter a setting value',
    text: 'Text', boolean: 'Boolean', integer: 'Integer', number: 'Number', json: 'JSON', note: 'Description',
    notePlaceholder: 'Help administrators understand this setting', noteHint: 'Use stable, readable keys; do not put sensitive credentials here.',
    create: 'Create setting', configured: 'Configured settings', configuredDescription: 'Showing settings returned by the Memory or GORM provider.',
    loading: 'Loading settings…', empty: 'No settings', emptyDescription: 'Create a namespace setting to get started.',
    missingNote: 'No description', saving: 'Saving…', save: 'Save', loadFailed: 'Settings could not be loaded. Please try again.',
    saveFailed: 'Settings could not be saved. Check the input.', createFailed: 'Setting could not be created. Check the input.', required: 'Enter both a namespace and key.',
    saved: 'Saved {id}', created: 'Setting created',
  },
}
