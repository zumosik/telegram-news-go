package bot

const (
	helpMsg = `
	Admin commands:
	/add NAME, URL - adding new source
	/list - list of all sources
	/delete ID - delete source by id
	`
	helloMsg         = "Hello! 👋\n" + helpMsg
	sourceIsAddedMsg = "👌 Source is added with ID: `%d`\\. You can use this ID to manage source\\."
	sorceIsDeleted   = "👌 Source id deleted\\."
	listSourcesMsg   = "List of sources (total %d):\n\n%s"
	invalidDataMsg   = "❌ Invalid data, check /help ❌"
)
