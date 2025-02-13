package owner

import "context"

type Command struct {
	r Repository
}

func NewCommand(r Repository) Command {
	return Command{r: r}
}

func (c Command) NewOwner(ctx context.Context, customerID string) error {
	_, err := c.r.FindOwnerByID(ctx, customerID)
	if err != nil {
		return err
	}

	isOwner, err := c.r.IsOwner(ctx, customerID)
	if err != nil {
		return err
	}
	if !isOwner {
		ow := NewOwner(customerID)
		err = c.r.SaveOwner(ctx, ow)
		if err != nil {
			return err
		}
	}
	return nil
}
