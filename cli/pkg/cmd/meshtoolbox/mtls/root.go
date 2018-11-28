package mtls

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/common"
	"github.com/solo-io/supergloo/cli/pkg/nsutil"
	superglooV1 "github.com/solo-io/supergloo/pkg/api/v1"
	"github.com/spf13/cobra"
)

// we have not yet implemented mtls updates
// if you are working on mTLS updates, set DEV_MTLS to true
// when you implement mTLS updates, remove the placeholder
// comment: consider passing this setting as a dev flag
// pros: easier to configure
// cons: may confuse users when shown in help msg, might not add a lot of value
const DEV_MTLS = false

// strings that users will pass to trigger commands
const (
	ENABLE_MTLS  = "enable"
	DISABLE_MTLS = "disable"
	TOGGLE_MTLS  = "toggle"
)

var validRootArgs = []string{ENABLE_MTLS, DISABLE_MTLS, TOGGLE_MTLS} // for bash completion

func Root(opts *options.Options) *cobra.Command {
	if !DEV_MTLS {
		cmd := &cobra.Command{
			Use:   "mtls",
			Short: `Set mTLS status`,
			Long:  `Set mTLS status`,
			RunE: func(c *cobra.Command, args []string) error {
				// this function does nothing but it triggers validation
				fmt.Println("Warning: mTLS config is not yet available. In the meantime, you can specify mTLS properties during install.")
				return nil
			},
		}
		return cmd
	}
	cmd := &cobra.Command{
		Use:       "mtls",
		Short:     `set mTLS status`,
		Long:      `set mTLS status`,
		ValidArgs: validRootArgs,
		Args:      rootArgValidation,
		RunE: func(c *cobra.Command, args []string) error {
			// this function does nothing but it triggers validation
			return nil
		},
	}
	cmd.AddCommand(
		Enable(opts),
		Disable(opts),
		Toggle(opts),
	)
	return cmd
}

func rootArgValidation(c *cobra.Command, args []string) error {
	expectedArgCount := 1
	if len(args) != expectedArgCount {
		return fmt.Errorf("Too many args (%v given, %v expected)", len(args), expectedArgCount)
	}
	subCommandName := args[0]
	if !common.Contains(validRootArgs, subCommandName) {
		return fmt.Errorf("%v is not a valid argument", subCommandName)
	}
	return nil
}

func Enable(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   ENABLE_MTLS,
		Short: `enable mTLS`,
		Long:  `enable mTLS`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := enableMtls(opts); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func Disable(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   DISABLE_MTLS,
		Short: `disable mTLS`,
		Long:  `disable mTLS`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := disableMtls(opts); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func Toggle(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   TOGGLE_MTLS,
		Short: `toggle mTLS`,
		Long:  `toggle mTLS`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := toggleMtls(opts); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func enableMtls(opts *options.Options) error {

	if _, err := updateMtls(ENABLE_MTLS, opts); err != nil {
		return err
	}
	fmt.Printf("Enabled mTLS on mesh %v", opts.MeshTool.Mesh.Name)

	return nil
}

func disableMtls(opts *options.Options) error {
	if _, err := updateMtls(DISABLE_MTLS, opts); err != nil {
		return err
	}
	fmt.Printf("Disabled mTLS on mesh %v", opts.MeshTool.Mesh.Name)
	return nil
}

func toggleMtls(opts *options.Options) error {
	mesh, err := updateMtls(TOGGLE_MTLS, opts)
	if err != nil {
		return err
	}
	status := "disabled"
	if mesh.Encryption.TlsEnabled {
		status = "enabled"
	}
	fmt.Printf("Toggled (%v) mTLS on mesh %v", status, opts.MeshTool.Mesh.Name)
	return nil
}

// Ensure that all the needed user-specified values have been provided
func ensureFlags(operation string, opts *options.Options) error {

	// all operations require a target mesh spec
	meshRef := &(opts.MeshTool).Mesh
	if err := nsutil.EnsureMesh(meshRef, opts); err != nil {
		return err
	}

	return nil
}

func updateMtls(operation string, opts *options.Options) (*superglooV1.Mesh, error) {
	// 1. validate/aquire arguments
	if err := ensureFlags(operation, opts); err != nil {
		return nil, err
	}

	// 2. read the existing mesh
	meshClient, err := common.GetMeshClient()
	if err != nil {
		return nil, err
	}
	meshRef := &(opts.MeshTool).Mesh
	mesh, err := (*meshClient).Read(meshRef.Namespace, meshRef.Name, clients.ReadOpts{})
	if err != nil {
		return nil, err
	}

	// 3. mutate the mesh structure
	switch operation {
	case ENABLE_MTLS:
		if mesh.Encryption == nil {
			mesh.Encryption = &superglooV1.Encryption{
				TlsEnabled: true,
			}
		} else {
			mesh.Encryption.TlsEnabled = true

		}
	case DISABLE_MTLS:
		if mesh.Encryption == nil {
			mesh.Encryption = &superglooV1.Encryption{
				TlsEnabled: false,
			}
		} else {
			mesh.Encryption.TlsEnabled = false

		}
	case TOGGLE_MTLS:
		// if encryption has not been specified, "toggle" will enable it
		if mesh.Encryption == nil {
			mesh.Encryption = &superglooV1.Encryption{
				TlsEnabled: true,
			}
		} else {
			mesh.Encryption.TlsEnabled = !mesh.Encryption.TlsEnabled

		}
	default:
		panic(fmt.Errorf("Operation %v not recognized", operation))
	}

	// 4. write the changes
	writtenMesh, err := (*meshClient).Write(mesh, clients.WriteOpts{OverwriteExisting: true})
	if err != nil {
		return nil, err
	}
	return writtenMesh, nil
}
