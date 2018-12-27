package pligos

type Infrastructure struct {
	ValuesEncoderMap  ValuesEncoderMap
	SecretProviderMap SecretProviderMap

	config *PligosConfig
	Friend *Friend
}

func NewInfrastructure(config *PligosConfig, friend *Friend, ve ValuesEncoderMap, sp SecretProviderMap) *Infrastructure {
	return &Infrastructure{
		ValuesEncoderMap:  ve,
		SecretProviderMap: sp,
		config:            config,
		Friend:            friend,
	}
}

func (infra *Infrastructure) Write() error {
	for _, context := range infra.config.Contexts {
		values, err := infra.ValuesEncoderMap.Get(context.Name).EncodeContext(context.Name)
		if err != nil {
			return err
		}

		values = infra.Friend.EnrichValues(values)

		helm, err := NewHelm(infra.SecretProviderMap.Get(context.Name), context.Flavor, context.Output, values)
		if err != nil {
			return err
		}

		defer helm.Clean()

		desc := HelmDescription{Description: infra.config.Description, Version: infra.config.ChartVersion}
		if err := helm.Create(desc); err != nil {
			return err
		}
	}

	return nil
}
