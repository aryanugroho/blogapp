package console

import (
	"context"
	"errors"
	"fmt"

	"github.com/aryanugroho/blogapp/config"
	"github.com/aryanugroho/blogapp/infrastructure/sqlstore"
	"gopkg.in/ukautz/clif.v1"
)

var conf = config.All()

func (c *Console) MigrateCreate() *clif.Command {
	return clif.NewCommand("migrate:create", "Create a new migration file.", func(o *clif.Command, in clif.Input, out clif.Output) error {
		name := o.Argument("name").String()

		if err := c.Migrator.Create(name, sqlstore.Migration); err != nil {
			return err
		}

		return nil
	}).NewArgument("name", "Migration name", "", true, false)
}

func (c *Console) MigrateRun(ctx context.Context) *clif.Command {
	return clif.NewCommand("migrate:run", "Run the database migrations.", func(o *clif.Command, in clif.Input, out clif.Output) error {
		if isNotInProduction() {
			return c.Migrator.MigrateAll(ctx)
		}

		if confirmInProduction(in) {
			return c.Migrator.MigrateAll(ctx)
		}

		out.Printf("Migrate aborted.\n")

		return nil
	})
}

func (c *Console) MigrateRollback() *clif.Command {
	return clif.NewCommand("migrate:rollback", "Rollback the last database migration.", func(o *clif.Command, in clif.Input, out clif.Output) error {
		step := o.Option("step").Int()

		if step <= 0 {
			return errors.New("step can't be zero or negative")
		}

		if isNotInProduction() {
			return c.Migrator.Rollback(step)
		}

		if confirmInProduction(in) {
			return c.Migrator.Rollback(step)
		}

		out.Printf("Rollback aborted.\n")

		return nil
	}).NewOption("step", "s", "The number of migrations to be reverted", "1", false, false)
}

func (c *Console) MigrateReset() *clif.Command {
	return clif.NewCommand("migrate:reset", "Rollback all database migrations.", func(o *clif.Command, in clif.Input, out clif.Output) error {
		if isNotInProduction() {
			return c.Migrator.Reset()
		}

		if confirmInProduction(in) {
			return c.Migrator.Reset()
		}

		out.Printf("Reset aborted.\n")

		return nil
	})
}

func (c *Console) MigrateRefresh() *clif.Command {
	return clif.NewCommand("migrate:refresh", "Reset and re-run all migrations.", func(o *clif.Command, in clif.Input, out clif.Output) error {
		if isNotInProduction() {
			return c.Migrator.Refresh()
		}

		if confirmInProduction(in) {
			return c.Migrator.Refresh()
		}

		out.Printf("Refresh aborted.\n")

		return nil
	})
}

func confirmInProduction(in clif.Input) bool {
	fmt.Println("**************************************")
	fmt.Println("*     Application In Production!     *")
	fmt.Println("**************************************")

	return in.Confirm("Do you really wish to run this command?")
}

func isNotInProduction() bool {
	return conf.App.ENV != "prod"
}
