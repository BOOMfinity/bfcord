package slash

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
	"slices"
)

type OptionBuilder interface {
	Build() Option
}

type baseBuilder[B any] interface {
	NameLocalization(lang string, name string) B
	NameLocalizations(langs map[string]string) B
	DescriptionLocalization(lang string, name string) B
	DescriptionLocalizations(langs map[string]string) B
}

type CommandBuilder baseCommandBuilder[CommandBuilder]

type EditCommandBuilder interface {
	baseCommandBuilder[EditCommandBuilder]
	Name(str string) EditCommandBuilder
	Description(str string) EditCommandBuilder
}

type baseCommandBuilder[B any] interface {
	baseBuilder[B]
	DM(enabled bool) B
	DefaultPermission(allowed bool) B
	DefaultMemberPermissions(perms permissions.Permission) B
	Option(bl OptionBuilder) B
	Import(cmd Command) B
	Run() (Command, error)
}

type baseOptionBuilder[B any] interface {
	baseBuilder[B]
	Required() B
	Import(opt Option) B
	Build() Option
}

type SubCommandBuilder interface {
	baseBuilder[SubCommandBuilder]
	Option(bl OptionBuilder) SubCommandBuilder
	Build() Option
}

type SubCommandGroupBuilder interface {
	baseBuilder[SubCommandGroupBuilder]
	Option(bl OptionBuilder) SubCommandGroupBuilder
	Build() Option
}

type StringBuilder[A any] interface {
	baseOptionBuilder[StringBuilder[A]]
	MinLength(len uint16) StringBuilder[A]
	MaxLength(len uint16) StringBuilder[A]
	AutoComplete() StringBuilder[A]
	Choices(list []Choice) StringBuilder[A]
}

type BooleanBuilder baseOptionBuilder[BooleanBuilder]

type UserBuilder baseOptionBuilder[UserBuilder]

type RoleBuilder baseOptionBuilder[RoleBuilder]

type MentionableBuilder baseOptionBuilder[MentionableBuilder]

type AttachmentBuilder baseOptionBuilder[AttachmentBuilder]

type NumberBuilder[A any] interface {
	baseOptionBuilder[NumberBuilder[A]]
	MaxValue(val A) NumberBuilder[A]
	MinValue(val A) NumberBuilder[A]
	AutoComplete() NumberBuilder[A]
	Choices(list []Choice) NumberBuilder[A]
}

type DoubleBuilder[A any] interface {
	baseOptionBuilder[DoubleBuilder[A]]
	MaxValue(val A) DoubleBuilder[A]
	MinValue(val A) DoubleBuilder[A]
	AutoComplete() DoubleBuilder[A]
	Choices(list []Choice) DoubleBuilder[A]
}

type ChannelBuilder interface {
	baseOptionBuilder[ChannelBuilder]
	Types(allowed ...discord.ChannelType) ChannelBuilder
}

type OptionSelector interface {
	String() StringBuilder[string]
	Number() NumberBuilder[int]
	Double() DoubleBuilder[float64]
	Channel() ChannelBuilder
	Boolean() BooleanBuilder
	User() UserBuilder
	Role() RoleBuilder
	Mentionable() MentionableBuilder
	Attachment() AttachmentBuilder
	SubCommand() SubCommandBuilder
	SubCommandGroup() SubCommandGroupBuilder
}

var _ = (CommandBuilder)(&commandBuilder[CommandBuilder]{})
var _ = (EditCommandBuilder)(&commandBuilder[EditCommandBuilder]{})

type commandBuilder[B any] struct {
	b      B
	cmd    *Command
	client *api.Client
	id     snowflake.ID
	guild  snowflake.ID
	app    snowflake.ID
}

func (c *commandBuilder[B]) Run() (cmd Command, err error) {
	req := c.client.New(true)
	url := fmt.Sprintf("%v/applications/%v", api.FullApiUrl, c.app)
	if c.guild.Valid() {
		url += fmt.Sprintf("/guilds/%v/commands", c.guild)
	} else {
		url += "/commands"
	}
	if c.id.Valid() {
		req.Header.SetMethod(fasthttp.MethodPatch)
		url += fmt.Sprintf("/%v", c.id)
	} else {
		req.Header.SetMethod(fasthttp.MethodPost)
	}
	body, err := json.Marshal(c.cmd)
	if err != nil {
		return cmd, fmt.Errorf("error while parsing command to json: %w", err)
	}
	req.SetBody(body)
	req.SetRequestURI(url)
	err = c.client.DoResult(req, &cmd)
	return
}

func (c *commandBuilder[B]) Name(str string) B {
	c.cmd.Name = str
	return c.b
}

func (c *commandBuilder[B]) Description(str string) B {
	c.cmd.Description = str
	return c.b
}

func (c *commandBuilder[B]) Import(cmd Command) B {
	*c.cmd = cmd
	return c.b
}

func (c *commandBuilder[B]) NameLocalization(lang string, name string) B {
	c.cmd.NameLocalizations[lang] = name
	return c.b
}

func (c *commandBuilder[B]) NameLocalizations(langs map[string]string) B {
	c.cmd.NameLocalizations = langs
	return c.b
}

func (c *commandBuilder[B]) DescriptionLocalization(lang string, name string) B {
	c.cmd.DescriptionLocalizations[lang] = name
	return c.b
}

func (c *commandBuilder[B]) DescriptionLocalizations(langs map[string]string) B {
	c.cmd.DescriptionLocalizations = langs
	return c.b
}

func (c *commandBuilder[B]) DM(enabled bool) B {
	c.cmd.DM = enabled
	return c.b
}

func (c *commandBuilder[B]) DefaultPermission(allowed bool) B {
	c.cmd.DefaultPermission = allowed
	return c.b
}

func (c *commandBuilder[B]) DefaultMemberPermissions(perms permissions.Permission) B {
	c.cmd.DefaultMemberPermissions = perms
	return c.b
}

func (c *commandBuilder[B]) Option(bl OptionBuilder) B {
	c.cmd.Options = append(c.cmd.Options, bl.Build())
	slices.SortStableFunc(c.cmd.Options, func(a, b Option) int {
		if a.Required && !b.Required {
			return 1
		}
		if b.Required && !a.Required {
			return -1
		}
		return 0
	})
	return c.b
}

var _ = (baseOptionBuilder[any])(&optionBuilder[any, any]{})
var _ = (StringBuilder[string])(&optionBuilder[StringBuilder[string], string]{})
var _ = (NumberBuilder[int])(&optionBuilder[NumberBuilder[int], int]{})
var _ = (DoubleBuilder[float64])(&optionBuilder[DoubleBuilder[float64], float64]{})
var _ = (MentionableBuilder)(&optionBuilder[MentionableBuilder, any]{})
var _ = (UserBuilder)(&optionBuilder[UserBuilder, any]{})
var _ = (AttachmentBuilder)(&optionBuilder[AttachmentBuilder, any]{})
var _ = (RoleBuilder)(&optionBuilder[RoleBuilder, any]{})
var _ = (BooleanBuilder)(&optionBuilder[BooleanBuilder, any]{})
var _ = (SubCommandBuilder)(&optionBuilder[SubCommandBuilder, any]{})
var _ = (SubCommandGroupBuilder)(&optionBuilder[SubCommandGroupBuilder, any]{})
var _ = (ChannelBuilder)(&optionBuilder[ChannelBuilder, any]{})

type optionBuilder[B any, A any] struct {
	option *Option
	b      B
}

func (o *optionBuilder[B, A]) Import(opt Option) B {
	*o.option = opt
	return o.b
}

func (o *optionBuilder[B, A]) Option(bl OptionBuilder) B {
	o.option.Options = append(o.option.Options, bl.Build())
	slices.SortStableFunc(o.option.Options, func(a, b Option) int {
		if a.Required && !b.Required {
			return 1
		}
		if b.Required && !a.Required {
			return -1
		}
		return 0
	})
	return o.b
}

func (o *optionBuilder[B, A]) Types(t ...discord.ChannelType) B {
	o.option.ChannelTypes = t
	return o.b
}

func (o *optionBuilder[B, A]) MaxValue(val A) B {
	o.option.MaxValue = val
	return o.b
}

func (o *optionBuilder[B, A]) MinValue(val A) B {
	o.option.MinValue = val
	return o.b
}

func (o *optionBuilder[B, A]) MinLength(len uint16) B {
	o.option.MinLength = len
	return o.b
}

func (o *optionBuilder[B, A]) MaxLength(len uint16) B {
	o.option.MaxLength = len
	return o.b
}

func (o *optionBuilder[B, A]) AutoComplete() B {
	o.option.Autocomplete = true
	return o.b
}

func (o *optionBuilder[B, A]) Choices(list []Choice) B {
	o.option.Choices = list
	return o.b
}

func (o *optionBuilder[B, A]) NameLocalization(lang string, name string) B {
	o.option.NameLocalizations[lang] = name
	return o.b
}

func (o *optionBuilder[B, A]) NameLocalizations(langs map[string]string) B {
	o.option.NameLocalizations = langs
	return o.b
}

func (o *optionBuilder[B, A]) DescriptionLocalization(lang string, name string) B {
	o.option.DescriptionLocalizations[lang] = name
	return o.b
}

func (o *optionBuilder[B, A]) DescriptionLocalizations(langs map[string]string) B {
	o.option.DescriptionLocalizations = langs
	return o.b
}

func (o *optionBuilder[B, A]) Required() B {
	o.option.Required = true
	return o.b
}

func (o *optionBuilder[B, A]) Build() Option {
	return *o.option
}

type selector struct {
	name, desc string
}

func (s selector) String() StringBuilder[string] {
	return newStringBuilder(s.name, s.desc)
}

func (s selector) Number() NumberBuilder[int] {
	return newNumberBuilder(s.name, s.desc)
}

func (s selector) Double() DoubleBuilder[float64] {
	return newDoubleBuilder(s.name, s.desc)
}

func (s selector) Channel() ChannelBuilder {
	return newChannelBuilder(s.name, s.desc)
}

func (s selector) Boolean() BooleanBuilder {
	return newBooleanBuilder(s.name, s.desc)
}

func (s selector) User() UserBuilder {
	return newUserBuilder(s.name, s.desc)
}

func (s selector) Role() RoleBuilder {
	return newRoleBuilder(s.name, s.desc)
}

func (s selector) SubCommandGroup() SubCommandGroupBuilder {
	return newSubCommandGroupBuilder(s.name, s.desc)
}

func (s selector) Mentionable() MentionableBuilder {
	return newMentionableBuilder(s.name, s.desc)
}

func (s selector) Attachment() AttachmentBuilder {
	return newAttachmentBuilder(s.name, s.desc)
}

func (s selector) SubCommand() SubCommandBuilder {
	return newSubCommandBuilder(s.name, s.desc)
}

func NewOption(name, desc string) selector {
	return selector{name, desc}
}

func newBuilder[B, A any](name, desc string) *optionBuilder[B, A] {
	opt := &optionBuilder[B, A]{}
	opt.option = new(Option)
	opt.option.Name = name
	opt.option.Description = desc
	return opt
}

func newStringBuilder(name, desc string) StringBuilder[string] {
	bl := newBuilder[StringBuilder[string], string](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeString
	return bl
}

func newNumberBuilder(name, desc string) NumberBuilder[int] {
	bl := newBuilder[NumberBuilder[int], int](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeInteger
	return bl
}

func newDoubleBuilder(name, desc string) DoubleBuilder[float64] {
	bl := newBuilder[DoubleBuilder[float64], float64](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeDouble
	return bl
}

func newUserBuilder(name, desc string) UserBuilder {
	bl := newBuilder[UserBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeUser
	return bl
}

func newBooleanBuilder(name, desc string) BooleanBuilder {
	bl := newBuilder[BooleanBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeBoolean
	return bl
}

func newChannelBuilder(name, desc string) ChannelBuilder {
	bl := newBuilder[ChannelBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeChannel
	return bl
}

func newMentionableBuilder(name, desc string) MentionableBuilder {
	bl := newBuilder[MentionableBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeMentionable
	return bl
}

func newRoleBuilder(name, desc string) RoleBuilder {
	bl := newBuilder[RoleBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeRole
	return bl
}

func newAttachmentBuilder(name, desc string) AttachmentBuilder {
	bl := newBuilder[AttachmentBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeAttachment
	return bl
}

func newSubCommandBuilder(name, desc string) SubCommandBuilder {
	bl := newBuilder[SubCommandBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeSubCommand
	return bl
}

func newSubCommandGroupBuilder(name, desc string) SubCommandGroupBuilder {
	bl := newBuilder[SubCommandGroupBuilder, any](name, desc)
	bl.b = bl
	bl.option.Type = OptionTypeSubCommandGroup
	return bl
}

func newCommandBuilder(cl *api.Client, app snowflake.ID, guild snowflake.ID, name, desc string) CommandBuilder {
	bl := &commandBuilder[CommandBuilder]{client: cl, guild: guild, app: app, cmd: new(Command)}
	bl.b = bl
	bl.cmd.Name = name
	bl.cmd.Description = desc
	bl.cmd.Type = CommandTypeChatInput
	return bl
}

func newEditCommandBuilder(cl *api.Client, app snowflake.ID, guild snowflake.ID, id snowflake.ID) EditCommandBuilder {
	bl := &commandBuilder[EditCommandBuilder]{id: id, guild: guild, app: app, client: cl, cmd: new(Command)}
	bl.b = bl
	return bl
}
