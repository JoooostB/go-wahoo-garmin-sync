package dropbox

import "time"

type Notification struct {
	ListFolder struct {
		Accounts []string `json:"accounts"`
	} `json:"list_folder"`
	Delta struct {
		Users []int `json:"users"`
	} `json:"delta"`
}

// {
// 	"list_folder": {
// 		"accounts": ["dbid:AAC5YDzltZ5Wb1xJUwIe7TmKdCHANgT5ZNo"]
// 	},
// 	"delta": {
// 		"users": [113703853]
// 	}
// }

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	UID          string `json:"uid"`
	AccountID    string `json:"account_id"`
}

type ListFolderRequest struct {
	Path                            string `json:"path"`
	Recursive                       bool   `json:"recursive"`
	IncludeMediaInfo                bool   `json:"include_media_info"`
	IncludeDeleted                  bool   `json:"include_deleted"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members"`
	IncludeMountedFolders           bool   `json:"include_mounted_folders"`
	IncludeNonDownloadableFiles     bool   `json:"include_non_downloadable_files"`
}

type ListFolderRequestContinue struct {
	Cursor string `json:"cursor"`
}

type ListFolderResponse struct {
	Entries []struct {
		Tag            string    `json:".tag"`
		Name           string    `json:"name"`
		ID             string    `json:"id"`
		ClientModified time.Time `json:"client_modified"`
		ServerModified time.Time `json:"server_modified"`
		Rev            string    `json:"rev"`
		Size           int       `json:"size"`
		PathLower      string    `json:"path_lower"`
		PathDisplay    string    `json:"path_display"`
		SharingInfo    struct {
			ReadOnly             bool   `json:"read_only"`
			ParentSharedFolderID string `json:"parent_shared_folder_id"`
			ModifiedBy           string `json:"modified_by"`
		} `json:"sharing_info"`
		IsDownloadable bool `json:"is_downloadable"`
		PropertyGroups []struct {
			TemplateID string `json:"template_id"`
			Fields     []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
		} `json:"property_groups"`
		HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
		ContentHash              string `json:"content_hash"`
		FileLockInfo             struct {
			IsLockholder   bool      `json:"is_lockholder"`
			LockholderName string    `json:"lockholder_name"`
			Created        time.Time `json:"created"`
		} `json:"file_lock_info"`
	} `json:"entries"`
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
}
