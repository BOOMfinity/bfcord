package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/BOOMfinity/bfcord/api/images"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

type lowLevelQuery struct {
	client *Client
	emptyOptions[discord.LowLevelClientQuery]
}

func (v lowLevelQuery) CreateForumMessage(id snowflake.ID, _data discord.ForumMessageCreate) (d *discord.ChannelWithMessage, err error) {
	req := v.client.New(true)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/threads", FullApiUrl, id))
	data := _data.Message
	if err != nil {
		return
	}
	if data.Files != nil && len(*data.Files) > 0 {
		if _err := v.prepareAttachments(req, data); _err != nil {
			err = fmt.Errorf("failed to prepare attachments: %w", _err)
			return
		}
	} else {
		rawData, _err := json.Marshal(_data)
		if _err != nil {
			err = fmt.Errorf("failed to marshal json data: %w", _err)
			return
		}
		fmt.Println(string(rawData))
		req.SetBody(rawData)
	}
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	err = v.client.DoResult(req, &d)
	if d != nil {
		d.Message.Patch()
	}
	return
}

func (v lowLevelQuery) CreateOrUpdate(guild, role snowflake.ID, data discord.RoleCreate) (r *discord.Role, err error) {
	req := v.client.New(true)
	if role.Valid() {
		req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/roles/%v", FullApiUrl, guild, role))
		req.Header.SetMethod(fasthttp.MethodPatch)
	} else {
		req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/roles", FullApiUrl, guild))
		req.Header.SetMethod(fasthttp.MethodPost)
	}
	json, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("failed to marshal json body: %w", err)
		return
	}
	req.SetBody(json)
	err = v.client.DoResult(req, &r)
	return
}

func (v lowLevelQuery) UpdateGuild(guild snowflake.ID, data discord.GuildUpdate) (g *discord.Guild, err error) {
	req := v.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v", FullApiUrl, guild))
	req.Header.SetMethod(fasthttp.MethodPatch)
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}
	req.SetBody(raw)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return g, v.client.DoResult(req, &g)
}

func (v lowLevelQuery) StartThread(channel snowflake.ID, message snowflake.ID, data discord.ThreadCreate) (ch *discord.Channel, err error) {
	req := v.client.New(true)
	if message.Valid() {
		req.SetRequestURI(fmt.Sprintf("%v/channels/%v/messages/%v/threads", FullApiUrl, channel, message))
	} else {
		req.SetRequestURI(fmt.Sprintf("%v/channels/%v/threads", FullApiUrl, channel))
	}
	req.Header.SetMethod(fasthttp.MethodPost)
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}
	req.SetBody(raw)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	err = v.client.DoResult(req, &ch)
	return
}

func (v lowLevelQuery) UpdateChannel(id snowflake.ID, data discord.ChannelUpdate) (ch *discord.Channel, err error) {
	req := v.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v", FullApiUrl, id))
	req.Header.SetMethod(fasthttp.MethodPatch)
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}
	req.SetBody(raw)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	err = v.client.DoResult(req, &ch)
	return
}

func (v lowLevelQuery) ExecuteWebhook(id snowflake.ID, token string, data discord.WebhookExecute, wait bool, thread snowflake.ID) (msg *discord.Message, err error) {
	params := url.Values{}
	params.Add("wait", strconv.FormatBool(wait))
	if thread.Valid() {
		params.Add("thread_id", thread.String())
	}
	return v.Message(fasthttp.MethodPost, fmt.Sprintf("%v/webhooks/%v/%v?%v", FullApiUrl, id, token, params.Encode()), data)
}

func (v lowLevelQuery) UpdateWebhookMessage(id snowflake.ID, token string, message snowflake.ID, data discord.MessageCreate, thread snowflake.ID) (msg *discord.Message, err error) {
	urlx := fmt.Sprintf("%v/webhooks/%v/%v/messages/%v", FullApiUrl, id, token, message)
	if thread.Valid() {
		params := url.Values{}
		params.Add("thread_id", thread.String())
		urlx += "?" + params.Encode()
	}
	return v.Message(fasthttp.MethodPatch, urlx, data)
}

func (v lowLevelQuery) UpdateMessage(channel snowflake.ID, message snowflake.ID, data discord.MessageCreate) (msg *discord.Message, err error) {
	return v.Message(fasthttp.MethodPatch, fmt.Sprintf("%v/channels/%v/messages/%v", FullApiUrl, channel, message), data)
}

func (v lowLevelQuery) prepareMultipart(writer *multipart.Writer, files *[]discord.MessageFile) error {
	for i := range *files {
		file := (*files)[i]
		fileWriter, err := writer.CreateFormFile(fmt.Sprintf("files[%v]", i+1), file.Name)
		if err != nil {
			return err
		}
		if file.Name == "" {
			file.Name = fmt.Sprintf("file-%v", i+1)
			if file.Url != "" {
				file.Name += filepath.Ext(file.Url)
			}
		}
		switch {
		case file.Reader != nil:
			rawData, err := io.ReadAll(file.Reader)
			if err != nil {
				return err
			}
			if _, err = fileWriter.Write(rawData); err != nil {
				return err
			}
		case file.Url != "":
			req := v.client.New(false)
			req.SetRequestURI(file.Url)
			if file.Base64 {
				img, err := images.New([]byte(file.Url))
				if err != nil {
					return err
				}
				if _, err = fileWriter.Write(img.Data); err != nil {
					return err
				}
			} else {
				rawData, err := v.client.DoBytes(req, WithRetries(0))
				if err != nil {
					return err
				}
				if _, err = fileWriter.Write(rawData); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v lowLevelQuery) prepareAttachments(req *fasthttp.Request, data discord.MessageCreate) error {
	attachments := make([]discord.Attachment, 0, len(*data.Files))
	for i := range *data.Files {
		file := (*data.Files)[i]
		attachments = append(attachments, discord.Attachment{
			ID:          snowflake.ID(i + 1),
			Ephemeral:   file.Ephemeral,
			Description: file.Description,
		})
	}
	if data.Attachments != nil {
		*data.Attachments = append(*data.Attachments, attachments...)
	} else {
		data.Attachments = &attachments
	}
	writer := multipart.NewWriter(req.BodyWriter())
	jsonWriter, err := writer.CreateFormField("payload_json")
	if err != nil {
		return err
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = jsonWriter.Write(rawData)
	if err != nil {
		return err
	}
	err = v.prepareMultipart(writer, data.Files)
	if err != nil {
		return err
	}
	if err = writer.Close(); err != nil {
		return err
	}
	req.Header.SetContentType(writer.FormDataContentType())
	return nil
}

func (v lowLevelQuery) Message(method string, url string, _data any) (msg *discord.Message, err error) {
	req := v.client.New(true)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	var data discord.MessageCreate
	switch _d := _data.(type) {
	case discord.MessageCreate:
		data = _d
	case discord.WebhookExecute:
		data = _d.MessageCreate
	}
	if err != nil {
		return
	}
	if data.Files != nil && len(*data.Files) > 0 {
		if _err := v.prepareAttachments(req, data); _err != nil {
			err = fmt.Errorf("failed to prepare attachments: %w", _err)
			return
		}
	} else {
		rawData, _err := json.Marshal(_data)
		if _err != nil {
			err = fmt.Errorf("failed to marshal json data: %w", _err)
			return
		}
		req.SetBody(rawData)
	}
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	err = v.client.DoResult(req, &msg)
	if msg != nil {
		msg.Patch()
	}
	return
}

func (v lowLevelQuery) CreateMessage(channel snowflake.ID, data discord.MessageCreate) (msg *discord.Message, err error) {
	return v.Message(fasthttp.MethodPost, fmt.Sprintf("%v/channels/%v/messages", FullApiUrl, channel), data)
}

func (v lowLevelQuery) CreateGuildChannel(guild snowflake.ID, data discord.ChannelUpdate) (ch *discord.Channel, err error) {
	req := v.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/channels", FullApiUrl, guild))
	req.Header.SetMethod(fasthttp.MethodPost)
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}
	req.SetBody(raw)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return ch, v.client.DoResult(req, &ch)
}

func (v lowLevelQuery) UpdateGuildMember(guild snowflake.ID, member snowflake.ID, data discord.MemberUpdate) (m *discord.MemberWithUser, err error) {
	req := v.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members/%v", FullApiUrl, guild, member))
	req.Header.SetMethod(fasthttp.MethodPatch)
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}
	req.SetBody(raw)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	err = v.client.DoResult(req, &m)
	if err != nil {
		return
	}
	if m != nil {
		m.GuildID = guild
		m.UserID = member
	}
	return
}

func (v lowLevelQuery) SendDM(channel snowflake.ID) discord.CreateMessageBuilder {
	return builders.NewCreateMessageBuilder(channel)
}
