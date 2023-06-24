package cloud

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/masquernya/nextcloud-notifications/config"
	"github.com/masquernya/nextcloud-notifications/storage"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Huge thanks to this tool: https://blog.kowalczyk.info/tools/xmltogo/
// without it, this would have been a nightmare

type GetAllTaskGroupsResponse struct {
	XMLName  xml.Name `xml:"multistatus"`
	Text     string   `xml:",chardata"`
	D        string   `xml:"d,attr"`
	S        string   `xml:"s,attr"`
	Cal      string   `xml:"cal,attr"`
	Cs       string   `xml:"cs,attr"`
	Oc       string   `xml:"oc,attr"`
	Nc       string   `xml:"nc,attr"`
	Response []struct {
		Text     string `xml:",chardata"`
		Href     string `xml:"href"`
		Propstat []struct {
			Text string `xml:",chardata"`
			Prop struct {
				Text         string `xml:",chardata"`
				Resourcetype struct {
					Text           string  `xml:",chardata"`
					Collection     string  `xml:"collection"`
					Calendar       *string `xml:"calendar"`
					ScheduleInbox  string  `xml:"schedule-inbox"`
					ScheduleOutbox string  `xml:"schedule-outbox"`
					TrashBin       struct {
						Text string `xml:",chardata"`
						X2   string `xml:"x2,attr"`
					} `xml:"trash-bin"`
				} `xml:"resourcetype"`
				Owner struct {
					Text string `xml:",chardata"`
					Href string `xml:"href"`
				} `xml:"owner"`
				CurrentUserPrivilegeSet struct {
					Text      string `xml:",chardata"`
					Privilege []struct {
						Text                        string `xml:",chardata"`
						Write                       string `xml:"write"`
						WriteProperties             string `xml:"write-properties"`
						WriteContent                string `xml:"write-content"`
						Unlock                      string `xml:"unlock"`
						Bind                        string `xml:"bind"`
						Unbind                      string `xml:"unbind"`
						WriteAcl                    string `xml:"write-acl"`
						Read                        string `xml:"read"`
						ReadAcl                     string `xml:"read-acl"`
						ReadCurrentUserPrivilegeSet string `xml:"read-current-user-privilege-set"`
						ReadFreeBusy                string `xml:"read-free-busy"`
						ScheduleDeliver             string `xml:"schedule-deliver"`
						ScheduleDeliverInvite       string `xml:"schedule-deliver-invite"`
						ScheduleDeliverReply        string `xml:"schedule-deliver-reply"`
						ScheduleQueryFreebusy       string `xml:"schedule-query-freebusy"`
						ScheduleSend                string `xml:"schedule-send"`
						ScheduleSendInvite          string `xml:"schedule-send-invite"`
						ScheduleSendReply           string `xml:"schedule-send-reply"`
						ScheduleSendFreebusy        string `xml:"schedule-send-freebusy"`
						SchedulePostVevent          string `xml:"schedule-post-vevent"`
						All                         string `xml:"all"`
					} `xml:"privilege"`
				} `xml:"current-user-privilege-set"`
				Getcontenttype      string `xml:"getcontenttype"`
				Getetag             string `xml:"getetag"`
				Displayname         string `xml:"displayname"`
				SyncToken           string `xml:"sync-token"`
				Invite              string `xml:"invite"`
				AllowedSharingModes struct {
					Text           string `xml:",chardata"`
					CanBeShared    string `xml:"can-be-shared"`
					CanBePublished string `xml:"can-be-published"`
				} `xml:"allowed-sharing-modes"`
				PublishURL    string `xml:"publish-url"`
				CalendarOrder struct {
					Text string `xml:",chardata"`
					X1   string `xml:"x1,attr"`
				} `xml:"calendar-order"`
				CalendarColor struct {
					Text string `xml:",chardata"`
					X1   string `xml:"x1,attr"`
				} `xml:"calendar-color"`
				Getctag                       string `xml:"getctag"`
				CalendarDescription           string `xml:"calendar-description"`
				CalendarTimezone              string `xml:"calendar-timezone"`
				SupportedCalendarComponentSet struct {
					Text string `xml:",chardata"`
					Comp struct {
						Text string `xml:",chardata"`
						Name string `xml:"name,attr"`
					} `xml:"comp"`
				} `xml:"supported-calendar-component-set"`
				SupportedCalendarData struct {
					Text         string `xml:",chardata"`
					CalendarData []struct {
						Text        string `xml:",chardata"`
						ContentType string `xml:"content-type,attr"`
						Version     string `xml:"version,attr"`
					} `xml:"calendar-data"`
				} `xml:"supported-calendar-data"`
				MaxResourceSize         string `xml:"max-resource-size"`
				MinDateTime             string `xml:"min-date-time"`
				MaxDateTime             string `xml:"max-date-time"`
				MaxInstances            string `xml:"max-instances"`
				MaxAttendeesPerInstance string `xml:"max-attendees-per-instance"`
				SupportedCollationSet   struct {
					Text               string   `xml:",chardata"`
					SupportedCollation []string `xml:"supported-collation"`
				} `xml:"supported-collation-set"`
				CalendarFreeBusySet    string `xml:"calendar-free-busy-set"`
				ScheduleCalendarTransp struct {
					Text   string `xml:",chardata"`
					Opaque string `xml:"opaque"`
				} `xml:"schedule-calendar-transp"`
				ScheduleDefaultCalendarURL string `xml:"schedule-default-calendar-URL"`
				CalendarEnabled            string `xml:"calendar-enabled"`
				OwnerDisplayname           struct {
					Text string `xml:",chardata"`
					X2   string `xml:"x2,attr"`
				} `xml:"owner-displayname"`
				TrashBinRetentionDuration struct {
					Text string `xml:",chardata"`
					X2   string `xml:"x2,attr"`
				} `xml:"trash-bin-retention-duration"`
				DeletedAt struct {
					Text string `xml:",chardata"`
					X2   string `xml:"x2,attr"`
				} `xml:"deleted-at"`
				Source      string `xml:"source"`
				Refreshrate struct {
					Text string `xml:",chardata"`
					X1   string `xml:"x1,attr"`
				} `xml:"refreshrate"`
				SubscribedStripTodos       string `xml:"subscribed-strip-todos"`
				SubscribedStripAlarms      string `xml:"subscribed-strip-alarms"`
				SubscribedStripAttachments string `xml:"subscribed-strip-attachments"`
				CalendarAvailability       string `xml:"calendar-availability"`
			} `xml:"prop"`
			Status string `xml:"status"`
		} `xml:"propstat"`
	} `xml:"response"`
}

type GetAllTasksResponse struct {
	XMLName  xml.Name `xml:"multistatus"`
	Text     string   `xml:",chardata"`
	D        string   `xml:"d,attr"`
	S        string   `xml:"s,attr"`
	Cal      string   `xml:"cal,attr"`
	Cs       string   `xml:"cs,attr"`
	Oc       string   `xml:"oc,attr"`
	Nc       string   `xml:"nc,attr"`
	Response []struct {
		Text     string `xml:",chardata"`
		Href     string `xml:"href"`
		Propstat []struct {
			Text string `xml:",chardata"`
			Prop struct {
				Text           string `xml:",chardata"`
				Getcontenttype string `xml:"getcontenttype"`
				Getetag        string `xml:"getetag"`
				Resourcetype   string `xml:"resourcetype"`
				Owner          struct {
					Text string `xml:",chardata"`
					Href string `xml:"href"`
				} `xml:"owner"`
				CurrentUserPrivilegeSet struct {
					Text      string `xml:",chardata"`
					Privilege []struct {
						Text                        string `xml:",chardata"`
						Write                       string `xml:"write"`
						WriteProperties             string `xml:"write-properties"`
						WriteContent                string `xml:"write-content"`
						Unlock                      string `xml:"unlock"`
						WriteAcl                    string `xml:"write-acl"`
						Read                        string `xml:"read"`
						ReadAcl                     string `xml:"read-acl"`
						ReadCurrentUserPrivilegeSet string `xml:"read-current-user-privilege-set"`
					} `xml:"privilege"`
				} `xml:"current-user-privilege-set"`
				CalendarData string `xml:"calendar-data"`
				Displayname  string `xml:"displayname"`
				SyncToken    string `xml:"sync-token"`
			} `xml:"prop"`
			Status string `xml:"status"`
		} `xml:"propstat"`
	} `xml:"response"`
}

type RawTodo struct {
	UID          string
	Created      string
	LastModified string
	Due          string
	DtStamp      string
	Summary      string
	RelatedTo    string
}

func mustParse(str string) time.Time {
	// 20230618T130047
	t, err := time.Parse("20060102T150405", str)
	if err != nil {
		panic(err)
	}
	// I don't know what timezone this is. Are all nextcloud instances the same TZ? is there an API to get the timezone?
	t = t.UTC()
	t = t.Add(time.Hour * 4)
	return t
}

func (r *RawTodo) ToToDo() Todo {
	var due *time.Time
	if r.Due != "" {
		dt := mustParse(r.Due)
		due = &dt
	}

	return Todo{
		UID:     r.UID,
		Summary: r.Summary,
		Due:     due,
		Created: mustParse(r.Created),
		Updated: mustParse(r.LastModified),
	}
}

type Todo struct {
	UID     string
	Summary string
	Due     *time.Time
	Created time.Time
	Updated time.Time
}

func parseTodos(b string) ([]RawTodo, error) {
	todos := make([]RawTodo, 0)
	var currentTodo *RawTodo

	for _, line := range strings.Split(b, "\n") {
		if strings.HasPrefix(line, "BEGIN:VTODO") {
			currentTodo = &RawTodo{}
		} else if strings.HasPrefix(line, "UID:") {
			currentTodo.UID = line[4:]
		} else if strings.HasPrefix(line, "CREATED:") {
			currentTodo.Created = line[8:]
		} else if strings.HasPrefix(line, "LAST-MODIFIED:") {
			currentTodo.LastModified = line[14:]
		} else if strings.HasPrefix(line, "DTSTAMP:") {
			currentTodo.DtStamp = line[9:]
		} else if strings.HasPrefix(line, "SUMMARY:") {
			currentTodo.Summary = line[8:]
		} else if strings.HasPrefix(line, "RELATED-TO:") {
			currentTodo.RelatedTo = line[11:]
		} else if strings.HasPrefix(line, "DUE:") {
			currentTodo.Due = line[4:]
		} else if strings.HasPrefix(line, "END:VTODO") {
			todos = append(todos, *currentTodo)
		} else if strings.HasPrefix(line, "END:VCALENDAR") {
			break
		}
	}

	return todos, nil
}

func debugPrint(a any) string {
	bits, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(bits)
}

type TodoDefinition struct {
	Url   string
	Label string
	Id    string
}

func (c *Cloud) GetAllTaskGroups() ([]TodoDefinition, error) {
	for {
		// TODO: clean up
		req, err := http.NewRequest("PROPFIND", config.Get().CloudUrl+"/remote.php/dav/calendars/"+storage.Get().LoginUsername+"/", strings.NewReader(`<x0:propfind xmlns:x0="DAV:"><x0:prop><x0:getcontenttype/><x0:getetag/><x0:resourcetype/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/><x4:invite xmlns:x4="http://owncloud.org/ns"/><x5:allowed-sharing-modes xmlns:x5="http://calendarserver.org/ns/"/><x5:publish-url xmlns:x5="http://calendarserver.org/ns/"/><x6:calendar-order xmlns:x6="http://apple.com/ns/ical/"/><x6:calendar-color xmlns:x6="http://apple.com/ns/ical/"/><x5:getctag xmlns:x5="http://calendarserver.org/ns/"/><x1:calendar-description xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-timezone xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-component-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-data xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-resource-size xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:min-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-instances xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-attendees-per-instance xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-collation-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-free-busy-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-calendar-transp xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-default-calendar-URL xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x4:calendar-enabled xmlns:x4="http://owncloud.org/ns"/><x3:owner-displayname xmlns:x3="http://nextcloud.com/ns"/><x3:trash-bin-retention-duration xmlns:x3="http://nextcloud.com/ns"/><x3:deleted-at xmlns:x3="http://nextcloud.com/ns"/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/><x4:invite xmlns:x4="http://owncloud.org/ns"/><x5:allowed-sharing-modes xmlns:x5="http://calendarserver.org/ns/"/><x5:publish-url xmlns:x5="http://calendarserver.org/ns/"/><x6:calendar-order xmlns:x6="http://apple.com/ns/ical/"/><x6:calendar-color xmlns:x6="http://apple.com/ns/ical/"/><x5:getctag xmlns:x5="http://calendarserver.org/ns/"/><x1:calendar-description xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-timezone xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-component-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-data xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-resource-size xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:min-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-instances xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-attendees-per-instance xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-collation-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-free-busy-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-calendar-transp xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-default-calendar-URL xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x4:calendar-enabled xmlns:x4="http://owncloud.org/ns"/><x3:owner-displayname xmlns:x3="http://nextcloud.com/ns"/><x3:trash-bin-retention-duration xmlns:x3="http://nextcloud.com/ns"/><x3:deleted-at xmlns:x3="http://nextcloud.com/ns"/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/><x4:invite xmlns:x4="http://owncloud.org/ns"/><x5:allowed-sharing-modes xmlns:x5="http://calendarserver.org/ns/"/><x5:publish-url xmlns:x5="http://calendarserver.org/ns/"/><x6:calendar-order xmlns:x6="http://apple.com/ns/ical/"/><x6:calendar-color xmlns:x6="http://apple.com/ns/ical/"/><x5:getctag xmlns:x5="http://calendarserver.org/ns/"/><x1:calendar-description xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-timezone xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-component-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-data xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-resource-size xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:min-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-instances xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-attendees-per-instance xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-collation-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-free-busy-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-calendar-transp xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-default-calendar-URL xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x4:calendar-enabled xmlns:x4="http://owncloud.org/ns"/><x3:owner-displayname xmlns:x3="http://nextcloud.com/ns"/><x3:trash-bin-retention-duration xmlns:x3="http://nextcloud.com/ns"/><x3:deleted-at xmlns:x3="http://nextcloud.com/ns"/><x5:source xmlns:x5="http://calendarserver.org/ns/"/><x6:refreshrate xmlns:x6="http://apple.com/ns/ical/"/><x5:subscribed-strip-todos xmlns:x5="http://calendarserver.org/ns/"/><x5:subscribed-strip-alarms xmlns:x5="http://calendarserver.org/ns/"/><x5:subscribed-strip-attachments xmlns:x5="http://calendarserver.org/ns/"/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/><x4:invite xmlns:x4="http://owncloud.org/ns"/><x5:allowed-sharing-modes xmlns:x5="http://calendarserver.org/ns/"/><x5:publish-url xmlns:x5="http://calendarserver.org/ns/"/><x6:calendar-order xmlns:x6="http://apple.com/ns/ical/"/><x6:calendar-color xmlns:x6="http://apple.com/ns/ical/"/><x5:getctag xmlns:x5="http://calendarserver.org/ns/"/><x1:calendar-description xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-timezone xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-component-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-calendar-data xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-resource-size xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:min-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-date-time xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-instances xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:max-attendees-per-instance xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:supported-collation-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:calendar-free-busy-set xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-calendar-transp xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x1:schedule-default-calendar-URL xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x4:calendar-enabled xmlns:x4="http://owncloud.org/ns"/><x3:owner-displayname xmlns:x3="http://nextcloud.com/ns"/><x3:trash-bin-retention-duration xmlns:x3="http://nextcloud.com/ns"/><x3:deleted-at xmlns:x3="http://nextcloud.com/ns"/><x1:calendar-availability xmlns:x1="urn:ietf:params:xml:ns:caldav"/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/></x0:prop></x0:propfind>`))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Basic "+c.GetBasicAuth())
		req.Header.Set("Content-Type", "application/xml; charset=utf-8")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Depth", "1")
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, NewDavError(nil, nil, err)
		}
		bits, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 207 {
			return nil, NewDavError(resp, bits, nil)
		}
		var result []GetAllTaskGroupsResponse
		err = xml.Unmarshal(bits, &result)
		if err != nil {
			return nil, NewDavError(resp, bits, err)
		}
		calendarData := make([]TodoDefinition, 0)
		for _, item := range result {
			for _, response := range item.Response {
				//log.Info("data", debugPrint(response.Href))
				url := response.Href
				isTodo := false
				label := ""
				for _, prop := range response.Propstat {
					if prop.Prop.SupportedCalendarComponentSet.Comp.Name == "VTODO" {
						log.Info("Found TODO", url)
						isTodo = true
					}
					if prop.Prop.Displayname != "" {
						label = prop.Prop.Displayname
					}
				}
				if isTodo {
					id := strings.ReplaceAll(url, "/remote.php/dav/calendars/"+storage.Get().LoginUsername+"/", "")
					id = id[0:strings.Index(id, "/")]
					calendarData = append(calendarData, TodoDefinition{
						Label: label,
						Url:   url,
						Id:    id,
					})
				}
			}
		}
		log.Info("response", calendarData)
		return calendarData, nil
	}
}

func (c *Cloud) GetAllTasks(url string) ([]Todo, error) {
	for {
		// curl -H "OCS-APIREQUEST: true" -u <admin>:<admin-PW> -X POST <https://<my-cloud-address.de>/ocs/v2.php/apps/admin_notifications/api/v1/notifications/<user-to-send-message-to> -d "shortMessage=<message>"
		req, err := http.NewRequest("REPORT", config.Get().CloudUrl+url, strings.NewReader(`<x1:calendar-query xmlns:x1="urn:ietf:params:xml:ns:caldav"><x0:prop xmlns:x0="DAV:"><x0:getcontenttype/><x0:getetag/><x0:resourcetype/><x0:displayname/><x0:owner/><x0:resourcetype/><x0:sync-token/><x0:current-user-privilege-set/><x0:getcontenttype/><x0:getetag/><x0:resourcetype/><x1:calendar-data/></x0:prop><x1:filter><x1:comp-filter name="VCALENDAR"><x1:comp-filter name="VTODO"><x1:prop-filter name="completed"><x1:is-not-defined/></x1:prop-filter></x1:comp-filter></x1:comp-filter></x1:filter></x1:calendar-query>`))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Basic "+c.GetBasicAuth())
		req.Header.Set("Content-Type", "application/xml; charset=utf-8")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Depth", "1")
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, NewDavError(nil, nil, err)
		}
		bits, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 207 {
			return nil, NewDavError(resp, bits, nil)
		}
		var result []GetAllTasksResponse
		err = xml.Unmarshal(bits, &result)
		if err != nil {
			return nil, NewDavError(resp, bits, err)
		}
		calendarData := make([]Todo, 0)
		for _, item := range result {
			for _, response := range item.Response {
				for _, prop := range response.Propstat {
					if prop.Status == "HTTP/1.1 200 OK" {
						//log.Info("status", prop.Status, "for", prop.Prop.CalendarData)
						todo, err := parseTodos(strings.TrimSpace(prop.Prop.CalendarData))
						if err != nil {
							log.Info("skip invalid todo", err)
							continue
						}
						for _, t := range todo {
							calendarData = append(calendarData, t.ToToDo())
						}
						continue
					}
					//log.Info("prop", prop.Prop)
				}
			}
		}
		return calendarData, nil
	}
}

func (c *Cloud) DoSendNotification(t Todo) bool {
	if t.Due == nil {
		return false
	}
	if t.Due.After(time.Now().UTC()) {
		return false
	}
	log.Info("can be sent:", t.Due.Format(time.RFC3339), "vs", time.Now().UTC().Format(time.RFC3339))
	return true
}

type DiscordWebhookRequest struct {
	Content string `json:"content"`
}

func (c *Cloud) SendDiscordNotification(summary string, group string, uid string) {
	if config.Get().DiscordWebhook == "" {
		return
	}
	url := config.Get().CloudUrl + "/apps/tasks/#/calendars/" + group + "/tasks/" + uid + ".ics"
	log.Info("url", url)
	content := ""
	if config.Get().DiscordPingId != "" {
		content = "<@" + config.Get().DiscordPingId + "> "
	}
	content += "Task [\"" + summary + "\"](" + url + ") is due."
	bits, _ := json.Marshal(DiscordWebhookRequest{
		Content: content,
	})
	for {
		req, err := http.NewRequest("POST", config.Get().DiscordWebhook, bytes.NewReader(bits))
		if err != nil {
			panic(err)
		}
		req.Header.Set("content-type", "application/json")
		resp, err := c.client.Do(req)
		if err != nil {
			log.Info("error sending discord notification", err)
			time.Sleep(5 * time.Second)
			continue
		}
		bits, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Info("error read body in sending discord notification", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Info("error sending discord notification", resp.StatusCode, string(bits))
			time.Sleep(5 * time.Second)
			continue
		}
		return
	}
}

func (c *Cloud) SendNotification(shortMessage string, longMessage string, destinationUser string) {
	for {
		req, err := http.NewRequest("POST", config.Get().CloudUrl+"/ocs/v2.php/apps/notifications/api/v2/admin_notifications/"+destinationUser, strings.NewReader("shortMessage="+url.QueryEscape(shortMessage)+"&longMessage="+url.QueryEscape(longMessage)))
		if err != nil {
			panic(err)
		}
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("OCS-APIREQUEST", "true")
		req.Header.Set("Authorization", "Basic "+c.GetBasicAuthForAdmin())
		resp, err := c.client.Do(req)
		if err != nil {
			panic(err)
		}
		bits, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			panic(NewDavError(resp, bits, nil))
		}
		return
	}
}

func (c *Cloud) relTime(t time.Time) string {
	now := time.Now()
	if t.After(now) {
		return "in " + t.Sub(now).String()
	}
	return t.Sub(now).String() + " ago"
}

func (c *Cloud) SendNotifications() {
	tasks, err := c.GetAllTaskGroups()
	if err != nil {
		panic(err)
	}
	log.Info("tasks", tasks)
	sent := 0
	total := 0
	for _, task := range tasks {
		for {
			taskList, err := c.GetAllTasks(task.Url)
			if err != nil {
				panic(err)
			}
			//log.Info("taskList", taskList)
			for _, taskEntry := range taskList {
				total++
				if c.DoSendNotification(taskEntry) && storage.TrySetNotified(taskEntry.UID) {
					log.Info("notify for", taskEntry)
					//c.SendNotification("Task \""+taskEntry.Summary+"\" is due", "Task \""+taskEntry.Summary+"\" is due as of "+c.relTime(*taskEntry.Due), storage.Get().LoginUsername)
					c.SendDiscordNotification(taskEntry.Summary, task.Id, taskEntry.UID)
					sent++
				}
			}
			break
		}
	}
	log.Info("sent", sent, "notifications of", total, "tasks")
}
