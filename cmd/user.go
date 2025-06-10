package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewUserCmd creates the user command
func NewUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage Harbor users",
		Long:  `Manage Harbor users including listing, creation and deletion.`,
	}

	cmd.AddCommand(newUserListCmd())
	cmd.AddCommand(newUserCreateCmd())
	cmd.AddCommand(newUserDeleteCmd())

	return cmd
}

func newUserListCmd() *cobra.Command {
	var (
		search   string
		page     int
		pageSize int
		detail   bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long:  `List Harbor users.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			userSvc := harbor.NewUserService(client)

			opts := &api.ListOptions{Page: page, PageSize: pageSize}
			var users []*api.User
			if search != "" {
				users, err = userSvc.Search(search, opts)
			} else {
				users, err = userSvc.List(opts)
			}
			if err != nil {
				return fmt.Errorf("failed to list users: %w", err)
			}

			if len(users) == 0 {
				output.Info("No users found")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(users)
			case "yaml":
				return output.YAML(users)
			default:
				table := output.Table()
				headers := []string{"ID", "USERNAME", "EMAIL", "ADMIN", "CREATED"}
				if detail {
					headers = append(headers, "REALNAME")
				}
				table.Append(headers)

				for _, u := range users {
					row := []string{
						strconv.Itoa(u.UserID),
						u.Username,
						u.Email,
						strconv.FormatBool(u.SysadminFlag),
						u.CreationTime.Format("2006-01-02"),
					}
					if detail {
						row = append(row, u.Realname)
					}
					table.Append(row)
				}

				table.Render()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&search, "search", "", "Search by username")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Page size")
	cmd.Flags().BoolVar(&detail, "detail", false, "Show detailed information")

	return cmd
}

func newUserCreateCmd() *cobra.Command {
	var (
		email    string
		realname string
		password string
		admin    bool
	)

	cmd := &cobra.Command{
		Use:   "create <username>",
		Short: "Create a new user",
		Long:  `Create a new Harbor user.`,
		Args:  requireArgs(1, "requires <username>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			if password == "" {
				prompt := promptui.Prompt{Label: "Password", Mask: '*'}
				p, err := prompt.Run()
				if err != nil {
					return err
				}
				password = p
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			userSvc := harbor.NewUserService(client)

			req := &api.UserReq{
				Username: username,
				Email:    email,
				Password: password,
				Realname: realname,
			}

			user, err := userSvc.Create(req)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			if admin && user != nil {
				if err := userSvc.SetAdmin(int64(user.UserID), true); err != nil {
					output.Warning("Failed to set admin flag: %v", err)
				}
				user.SysadminFlag = true
			}

			output.Success("User '%s' created", username)
			if output.GetFormat() == "table" && user != nil {
				output.Info("")
				fmt.Printf("ID:        %d\n", user.UserID)
				fmt.Printf("Username:  %s\n", user.Username)
				fmt.Printf("Admin:     %v\n", user.SysadminFlag)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Email address")
	cmd.Flags().StringVar(&realname, "realname", "", "Real name")
	cmd.Flags().StringVar(&password, "password", "", "User password")
	cmd.Flags().BoolVar(&admin, "admin", false, "Create as Harbor admin")

	return cmd
}

func newUserDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <username>",
		Short: "Delete a user",
		Long:  `Delete a Harbor user by username.`,
		Args:  requireArgs(1, "requires <username>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			userSvc := harbor.NewUserService(client)

			user, err := userSvc.GetByUsername(username)
			if err != nil {
				return err
			}

			if !force {
				prompt := promptui.Prompt{Label: fmt.Sprintf("Delete user '%s'", username), IsConfirm: true}
				result, err := prompt.Run()
				if err != nil || strings.ToLower(result) != "y" {
					output.Info("Deletion cancelled")
					return nil
				}
			}

			if err := userSvc.Delete(int64(user.UserID)); err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}

			output.Success("User '%s' deleted", username)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force deletion without confirmation")
	return cmd
}
