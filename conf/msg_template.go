package conf

type MsgTemplate struct {
	ContractAuditSuccessAdvisor  int64 `yaml:"contract_audit_success_advisor"`
	ContractAuditSuccessUser     int64 `yaml:"contract_audit_success_user"`
	ContractAuditFailAdvisor     int64 `yaml:"contract_audit_fail_advisor"`
	ContractAuditFailUser        int64 `yaml:"contract_audit_fail_user"`
	QualifiedInvestorFailAdvisor int64 `yaml:"qualified_investor_fail_advisor"`
}
