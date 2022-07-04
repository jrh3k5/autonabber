# autonabber
Autonabber is a command-line tool that uses YAML files to apply budget changes in YNAB.

I created this tool because I do the same applications when my paychecks land and have been using a spreadsheet to track those applications. This manual process is, like any manual process, tedious and prone to error. By codifying the distributions of my paycheck - which rarely, if ever, change - into a configuration file, I can make sure that I am consistently applying my paychecks to the appropriate budget categories.

## Usage

To use this application, you will need the following:

* A personal access token (described [here](https://api.youneedabudget.com/#personal-access-tokens) in YNAB's documentatio)
* A YAML file describing the changes to be applied (refer to [YAML File Definition](#yaml-file-definition) for details on how to structure it)

Once you have these two artifacts, you can compile the application and then execute it so like so:

```
./autonabber --access-token=<personal access token> --file=<name of YAML file>
```

When you execute this application, you will be given the following prompts:

* If you have multiple budgets, you will be asked to select one
  * If you only have one budget, then you will not receive this prompt
* If you have multiple change sets in your YAML file, you will be asked to select one
  * If you only have one change set, then you will not receive this prompt
* You will be given a preview of your changes and asked to confirm
  * If you choose to confirm and the amount to be budgeted exceeds your funds in Ready to Assign, you will be prompted to confirm that you still wish to apply these changes
* If you have not opted to cancel the application of changes at any time, they will be applied to your budget

### Dry Run

To run this through to completion _except_ for the actual application of changes to the budget, specify the `-dry-run=true` option. This will still _read_ information from YNAB, but, if you confirm the application, it will not _write_ any changes to YNAB.

### Printing the Budget

If you want to see a copy of the budget as it's stored in YNAB, you can use the `-print-budget=true` option.

By default, hidden categories are not printed. If you want to see them, you can add the `-print-hidden-categories` option.

### YAML File Definition

Conceptually, the YAML file describes sets of changes to be applied - e.g., perhaps you get paid on the 1st and 15th day of each month and the way you distribute each paycheck amongst your budget categories between the 1st and the 15th. In this case, you would have two change sets: one to be applied on the 1st day of the month and one to be applied on the 15th day of the month.

(if you always apply the same changes, regardless of when the paycheck lands, then you would merely have one changeset to apply multiple times)

Refer to the [example](./example.yaml) for an example of what your YAML file can look like.

The structure of the YAML file is:

```yaml
changes
  - name: <name of change set>
    category_groups:
      - name: <name of category group as it appears in your budget in YNAB>
        categories:
          - name: <name of cateogry as it appears under the group in YNAB>
            change: <change operation>
```

The change operation can be one of the following:

* **Addition**: you can specify a change operation of `+##.##` to indicate by how much the budgeted balance for the category should be increased
* **Addition of Average Spent**: you can specify a change operation of `+average-spent-#m` to apply the monthly average spent (rounding up) in that category as an addition to the budgeted amount

## Development

This application has the following prerequisites:

* Go >= 1.18.0
* [gomock](https://github.com/golang/mock)

To build the application, run:

```
make build
```

To run the tests, run:

```
make test
```

To compile artifacts for release, run:

```
make release
```