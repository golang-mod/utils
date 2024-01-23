package conf

type SmsTemplate struct {
	ContractAuditFailAdvisor     int64 `yaml:"contract_audit_fail_advisor"`
	ContractAuditFailUser        int64 `yaml:"contract_audit_fail_user"`
	QualifiedInvestorFailAdvisor int64 `yaml:"qualified_investor_fail_advisor"`
	QualifiedInvestorFailUser    int64 `yaml:"qualified_investor_fail_user"`
	QualifiedInvestorSuccessUser int64 `yaml:"qualified_investor_success_user"`
}
