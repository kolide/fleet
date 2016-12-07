package mysql

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewAppConfig(info *kolide.AppConfig) (*kolide.AppConfig, error) {
	insertStatement := `
		INSERT INTO app_configs (
			org_name,
			org_logo_url,
			kolide_server_url,
			smtp_configured,
			smtp_sender_address,
			smtp_server,
			smtp_port,
			smtp_authentication_type,
			smtp_enable_ssl_tls,
			smtp_authentication_method,
			smtp_domain,
			smtp_user_name,
			smtp_password,
			smtp_verify_ssl_certs,
			smtp_enable_start_tls
		)
		VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )
	`

	err := d.db.Get(info, "SELECT * FROM app_configs LIMIT 1")
	switch err {
	case sql.ErrNoRows:
		result, err := d.db.Exec(insertStatement,
			info.OrgName,
			info.OrgLogoURL,
			info.KolideServerURL,
			info.Configured,
			info.SenderAddress,
			info.Server,
			info.Port,
			info.AuthenticationType,
			info.EnableSSLTLS,
			info.AuthenticationMethod,
			info.Domain,
			info.UserName,
			info.Password,
			info.VerifySSLCerts,
			info.EnableSSLTLS,
		)
		if err != nil {
			return nil, err
		}

		info.ID, _ = result.LastInsertId()
		return info, nil
	case nil:
		return info, d.SaveAppConfig(info)
	default:
		return nil, err
	}
}

func (d *Datastore) AppConfig() (*kolide.AppConfig, error) {
	info := &kolide.AppConfig{}
	err := d.db.Get(info, "SELECT * FROM app_configs LIMIT 1")
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (d *Datastore) SaveAppConfig(info *kolide.AppConfig) error {
	// only one row in table so no where clause
	sqlStatement := `
		UPDATE app_configs
		SET
			org_name = ?,
			org_logo_url = ?,
			kolide_server_url = ?,
			smtp_configured = ?,
			smtp_sender_address = ?,
			smtp_server = ?,
			smtp_port = ?,
			smtp_authentication_type = ?,
			smtp_enable_ssl_tls = ?,
			smtp_authentication_method = ?,
			smtp_domain = ?,
			smtp_user_name = ?,
			smtp_password = ?,
			smtp_verify_ssl_certs = ?,
			smtp_enable_start_tls = ?
	`
	_, err := d.db.Exec(sqlStatement,
		info.OrgName,
		info.OrgLogoURL,
		info.KolideServerURL,
		info.Configured,
		info.SenderAddress,
		info.Server,
		info.Port,
		info.AuthenticationType,
		info.EnableSSLTLS,
		info.AuthenticationMethod,
		info.Domain,
		info.UserName,
		info.Password,
		info.VerifySSLCerts,
		info.EnableStartTLS,
	)
	return err
}
