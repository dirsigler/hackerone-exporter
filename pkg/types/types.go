// Copyright 2025 Dennis Irsigler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import "time"

type Assets struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			AssetType                  string    `json:"asset_type"`
			DomainName                 string    `json:"domain_name"`
			Description                any       `json:"description"`
			Coverage                   string    `json:"coverage"`
			MaxSeverity                string    `json:"max_severity"`
			ConfidentialityRequirement string    `json:"confidentiality_requirement"`
			IntegrityRequirement       string    `json:"integrity_requirement"`
			AvailabilityRequirement    string    `json:"availability_requirement"`
			CreatedAt                  time.Time `json:"created_at"`
			UpdatedAt                  time.Time `json:"updated_at"`
			ArchivedAt                 time.Time `json:"archived_at"`
			Reference                  string    `json:"reference"`
			State                      string    `json:"state"`
		} `json:"attributes"`
		Relationships struct {
			AssetTags struct {
				AssetTagsData []struct {
					ID                      string `json:"id"`
					Type                    string `json:"type"`
					AssetTagsDataAttributes struct {
						Name string `json:"name"`
					} `json:"attributes"`
					AssetTagsDataRelationships struct {
						AssetTagCategory struct {
							AssetTagCategoryData struct {
								ID                             string `json:"id"`
								Type                           string `json:"type"`
								AssetTagCategoryDataAttributes struct {
									Name string `json:"name"`
								} `json:"attributes"`
							} `json:"data"`
						} `json:"asset_tag_category"`
					} `json:"relationships"`
				} `json:"data"`
			} `json:"asset_tags"`
			Programs struct {
				ProgramsData []struct {
					ID                     string `json:"id"`
					Type                   string `json:"type"`
					ProgramsDataAttributes struct {
						Handle string `json:"handle"`
						Name   string `json:"name"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"programs"`
			Attachments struct {
				AttachmentsData []struct {
					ID                        string `json:"id"`
					Type                      string `json:"type"`
					AttachmentsDataAttributes struct {
						ExpiringURL string    `json:"expiring_url"`
						CreatedAt   time.Time `json:"created_at"`
						FileName    string    `json:"file_name"`
						ContentType string    `json:"content_type"`
						FileSize    int       `json:"file_size"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"attachments"`
		} `json:"relationships"`
	} `json:"data"`
	Links struct{} `json:"links"`
}

type Reports struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title                    string    `json:"title"`
			State                    string    `json:"state"`
			CreatedAt                time.Time `json:"created_at"`
			SubmittedAt              time.Time `json:"submitted_at"`
			VulnerabilityInformation string    `json:"vulnerability_information"`
			TriagedAt                any       `json:"triaged_at"`
			ClosedAt                 any       `json:"closed_at"`
			LastReporterActivityAt   any       `json:"last_reporter_activity_at"`
			FirstProgramActivityAt   any       `json:"first_program_activity_at"`
			LastProgramActivityAt    any       `json:"last_program_activity_at"`
			BountyAwardedAt          any       `json:"bounty_awarded_at"`
			LastActivityAt           any       `json:"last_activity_at"`
			LastPublicActivityAt     any       `json:"last_public_activity_at"`
			SwagAwardedAt            any       `json:"swag_awarded_at"`
			DisclosedAt              any       `json:"disclosed_at"`
		} `json:"attributes,omitempty"`
		Relationships struct {
			Reporter struct {
				Data struct {
					ID         string `json:"id"`
					Type       string `json:"type"`
					Attributes struct {
						Username       string    `json:"username"`
						Name           string    `json:"name"`
						Disabled       bool      `json:"disabled"`
						CreatedAt      time.Time `json:"created_at"`
						ProfilePicture struct {
							Six2X62   string `json:"62x62"`
							Eight2X82 string `json:"82x82"`
							One10X110 string `json:"110x110"`
							Two60X260 string `json:"260x260"`
						} `json:"profile_picture"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"reporter"`
			Collaborators struct {
				Data []struct {
					Weight int `json:"weight"`
					User   struct {
						ID         string `json:"id"`
						Type       string `json:"type"`
						Attributes struct {
							Username       string    `json:"username"`
							Name           string    `json:"name"`
							Disabled       bool      `json:"disabled"`
							CreatedAt      time.Time `json:"created_at"`
							ProfilePicture struct {
								Six2X62   string `json:"62x62"`
								Eight2X82 string `json:"82x82"`
								One10X110 string `json:"110x110"`
								Two60X260 string `json:"260x260"`
							} `json:"profile_picture"`
							Reputation int `json:"reputation"`
							Signal     int `json:"signal"`
							Impact     int `json:"impact"`
						} `json:"attributes"`
					} `json:"user"`
				} `json:"data"`
			} `json:"collaborators"`
			Program struct {
				Data struct {
					ID         string `json:"id"`
					Type       string `json:"type"`
					Attributes struct {
						Handle    string    `json:"handle"`
						CreatedAt time.Time `json:"created_at"`
						UpdatedAt time.Time `json:"updated_at"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"program"`
			Weakness struct {
				Data struct {
					ID         string `json:"id"`
					Type       string `json:"type"`
					Attributes struct {
						Name        string    `json:"name"`
						Description string    `json:"description"`
						ExternalID  string    `json:"external_id"`
						CreatedAt   time.Time `json:"created_at"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"weakness"`
			Bounties struct {
				Data []any `json:"data"`
			} `json:"bounties"`
		} `json:"relationships,omitempty"`
	} `json:"data"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
		Last string `json:"last"`
	} `json:"links"`
}

type Programs struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Handle    string    `json:"handle"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"attributes"`
	} `json:"data"`
	Links struct{} `json:"links"`
}

type InvitedHackers struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			State       string    `json:"state"`
			CreatedAt   time.Time `json:"created_at"`
			ViewedAt    time.Time `json:"viewed_at"`
			AcceptedAt  time.Time `json:"accepted_at"`
			ExpiresAt   any       `json:"expires_at"`
			UpdatedAt   time.Time `json:"updated_at"`
			RejectedAt  any       `json:"rejected_at"`
			CancelledAt any       `json:"cancelled_at"`
		} `json:"attributes"`
		Relationships struct {
			Recipient struct {
				ID                  string `json:"id"`
				Type                string `json:"type"`
				RecipientAttributes struct {
					Username       string    `json:"username"`
					Name           string    `json:"name"`
					Disabled       bool      `json:"disabled"`
					CreatedAt      time.Time `json:"created_at"`
					ProfilePicture struct {
						Six2X62   string `json:"62x62"`
						Eight2X82 string `json:"82x82"`
						One10X110 string `json:"110x110"`
						Two60X260 string `json:"260x260"`
					} `json:"profile_picture"`
					Signal           any    `json:"signal"`
					Impact           any    `json:"impact"`
					Reputation       any    `json:"reputation"`
					Bio              string `json:"bio"`
					Website          string `json:"website"`
					Location         string `json:"location"`
					HackeroneTriager bool   `json:"hackerone_triager"`
				} `json:"attributes"`
			} `json:"recipient"`
			InvitedBy struct {
				ID                  string `json:"id"`
				Type                string `json:"type"`
				InvitedByAttributes struct {
					Username                          string    `json:"username"`
					Name                              string    `json:"name"`
					Disabled                          bool      `json:"disabled"`
					CreatedAt                         time.Time `json:"created_at"`
					InvitedByAttributesProfilePicture struct {
						Six2X62   string `json:"62x62"`
						Eight2X82 string `json:"82x82"`
						One10X110 string `json:"110x110"`
						Two60X260 string `json:"260x260"`
					} `json:"profile_picture"`
					Signal           int    `json:"signal"`
					Impact           int    `json:"impact"`
					Reputation       int    `json:"reputation"`
					Bio              string `json:"bio"`
					Website          string `json:"website"`
					Location         string `json:"location"`
					HackeroneTriager bool   `json:"hackerone_triager"`
				} `json:"attributes"`
			} `json:"invited_by"`
		} `json:"relationships"`
	} `json:"data"`
	Links struct{} `json:"links"`
}

type Weaknesses struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name        string    `json:"name"`
			Description string    `json:"description"`
			CreatedAt   time.Time `json:"created_at"`
			ExternalID  string    `json:"external_id"`
		} `json:"attributes"`
	} `json:"data"`
	Links struct{} `json:"links"`
}
