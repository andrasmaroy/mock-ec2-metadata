package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gizmo/web"
)

type (
	SecurityCredentials struct {
		User            string `json:"User"`
		AccessKeyId     string `json:"AccessKeyId"`
		SecretAccessKey string `json:"SecretAccessKey"`
		Token           string `json:"Token"`
		Expiration      string `json:"Expiration"`
		Code            string `json:"Code"`
	}

	MetadataValues struct {
		AmiId               string              `json:"ami-id"`
		AmiLaunchIndex      string              `json:"ami-launch-index"`
		AmiManifestPath     string              `json:"ami-manifest-path"`
		AvailabilityZone    string              `json:"availability-zone"`
		Hostname            string              `json:"hostname"`
		InstanceAction      string              `json:"instance-action"`
		InstanceId          string              `json:"instance-id"`
		InstanceType        string              `json:"instance-type"`
		LocalHostName       string              `json:"local-hostname"`
		LocalIpv4           string              `json:"local-ipv4"`
		Mac                 string              `json:"mac"`
		Profile             string              `json:"profile"`
		ReservationId       string              `json:"reservation-id"`
		SecurityGroups      []string            `json:"security-groups"`
		SecurityCredentials SecurityCredentials `json:"security-credentials"`
	}

	Config struct {
		Server           *server.Config
		MetadataValues   *MetadataValues
		MetadataPrefixes []string
		UserdataValues   map[string]string
		UserdataPrefixes []string
	}

	MetadataService struct {
		config *Config
	}
)

func NewMetadataService(cfg *Config) *MetadataService {
	return &MetadataService{cfg}
}

func (s *MetadataService) Middleware(h http.Handler) http.Handler {
	return h
}

// middleware for adding plaintext content type
func plainText(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		h(w, r)
	}
}

func movedPermanently(redirectPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
	}
}

func (s *MetadataService) GetAmiId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiId)
}

func (s *MetadataService) GetAmiLaunchIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiLaunchIndex)
}

func (s *MetadataService) GetAmiManifestPath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiManifestPath)
}

func (s *MetadataService) GetAvailabilityZone(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AvailabilityZone)
}

func (s *MetadataService) GetHostName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Hostname)
}

func (s *MetadataService) GetInstanceAction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceAction)
}

func (s *MetadataService) GetInstanceId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceId)
}

func (s *MetadataService) GetInstanceType(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceType)
}

func (s *MetadataService) GetLocalHostName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.LocalHostName)
}

func (s *MetadataService) GetLocalIpv4(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.LocalIpv4)
}

func (s *MetadataService) GetIAM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "security-credentials/")
}

func (s *MetadataService) GetMac(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Mac)
}

func (s *MetadataService) GetProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Profile)
}

func (s *MetadataService) GetReservationId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.ReservationId)
}

func (s *MetadataService) GetSecurityCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.SecurityCredentials.User)
}

func (s *MetadataService) GetSecurityGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, strings.Join(s.config.MetadataValues.SecurityGroups, "\n"))
}

func (s *MetadataService) GetSecurityCredentialDetails(w http.ResponseWriter, r *http.Request) {
	username := web.Vars(r)["username"]

	if username != s.config.MetadataValues.SecurityCredentials.User {
		server.Log.Error("error, IAM user not found")
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(s.config.MetadataValues.SecurityCredentials)
	if err != nil {
		server.Log.Error("error converting security credentails to json: ", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	server.LogWithFields(r).Info("GetSecurityCredentialDetails returning: %#v",
		s.config.MetadataValues.SecurityCredentials)
}

func (s *MetadataService) GetMetadata(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `meta-data`)
}

func (s *MetadataService) GetMetadataIndex(w http.ResponseWriter, r *http.Request) {
	var availableItems = []string{"ami-id", "ami-launch-index", "ami-manifest-path", "block-device-mapping/",
		"hostname", "iam/", "instance-action", "instance-id", "instance-type", "local-hostname",
		"local-ipv4", "mac", "metrics/", "network/", "placement/", "profile", "public-hostname",
		"public-ipv4", "public-keys/", "reservation-id", "security-groups", "services/"}

	fmt.Fprintf(w, strings.Join(availableItems, "\n"))
}

func (s *MetadataService) GetUserData(w http.ResponseWriter, r *http.Request) {

	for index, value := range s.config.UserdataValues {
		fmt.Fprintf(w, fmt.Sprint(index+"="+value+"\n"))
	}
}

func (s *MetadataService) GetIndex(w http.ResponseWriter, r *http.Request) {
	var metadataVersions []string
	for _, metadataPrefix := range s.config.MetadataPrefixes{
		metadataVersions = append(metadataVersions, strings.Split(metadataPrefix, "/")[1])
	}

	fmt.Fprintf(w, strings.Join(metadataVersions, "\n"))
}

// Endpoints is a listing of all endpoints available in the MetadataService.
func (service *MetadataService) Endpoints() map[string]map[string]http.HandlerFunc {
	handlers := map[string]map[string]http.HandlerFunc{}

	for index, metadataPrefix := range service.config.MetadataPrefixes {
		server.Log.Info("adding Metadata prefix (", index, ") ", metadataPrefix)

		var metadataVersion = strings.Split(metadataPrefix, "/")[1]
		server.Log.Info("adding metadata version: ", metadataVersion)
		handlers["/" + metadataVersion + "/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetMetadata),
		}
		handlers[metadataPrefix+"/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetMetadataIndex),
		}
		handlers[metadataPrefix+"/ami-id"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiId),
		}
		handlers[metadataPrefix+"/ami-launch-index"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiLaunchIndex),
		}
		handlers[metadataPrefix+"/ami-manifest-path"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiManifestPath),
		}
		handlers[metadataPrefix+"/placement/availability-zone"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAvailabilityZone),
		}
		handlers[metadataPrefix+"/hostname"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetHostName),
		}
		handlers[metadataPrefix+"/instance-action"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceAction),
		}
		handlers[metadataPrefix+"/instance-id"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceId),
		}
		handlers[metadataPrefix+"/instance-type"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceType),
		}
		handlers[metadataPrefix+"/iam/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetIAM),
		}
		handlers[metadataPrefix+"/iam/security-credentials"] = map[string]http.HandlerFunc{
			"GET": movedPermanently(metadataPrefix + "/iam/security-credentials/"),
		}
		handlers[metadataPrefix+"/iam/security-credentials/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityCredentials),
		}
		handlers[metadataPrefix+"/iam/security-credentials/{username}"] = map[string]http.HandlerFunc{
			"GET": service.GetSecurityCredentialDetails,
		}
		handlers[metadataPrefix+"/local-hostname"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetLocalHostName),
		}
		handlers[metadataPrefix+"/local-ipv4"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetLocalIpv4),
		}
		handlers[metadataPrefix+"/mac"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetMac),
		}
		handlers[metadataPrefix+"/profile"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetProfile),
		}
		handlers[metadataPrefix+"/reservation-id"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetReservationId),
		}
		handlers[metadataPrefix+"/security-groups"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityGroups),
		}
	}

	for index, value := range service.config.UserdataPrefixes {
		server.Log.Info("adding Userdata prefix (", index, ") ", value)

		handlers[value+"/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetUserData),
		}
	}
	handlers["/"] = map[string]http.HandlerFunc{
		"GET": service.GetIndex,
	}
	return handlers
}

func (s *MetadataService) Prefix() string {
	return "/"
}

type error struct {
	Err string
}

func (e *error) Error() string {
	return e.Err
}
