/*
alertlogic provides an API to work with the Alert Logic Cloud Insights API.
This is currently in _very early_ development and only provides minimal available endpoints.

Create a new API client:

	api, err := alertlogic.NewWithUsernameAndPassword(
		os.Getenv("ALERTLOGIC_ACCOUNT_ID"),
		os.Getenv("ALERTLOGIC_USERNAME"),
		os.Getenv("ALERTLOGIC_PASSWORD"),
	)


Get account details:

	resp, err := api.GetAccountDetails()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)
*/
package alertlogic
