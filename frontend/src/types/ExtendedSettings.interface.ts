export type SettingInputType = 'FolderPicker' | 'APIPicker' | 'APIModelPicker' | 'ExtendedTextInput' | 'NumberInput' | 'Boolean' | 'URLInput' | 'TimePicker'

export interface ExtendedSettingDisplayProps {
  DisplayName: string
  Description: string
  Category: string
  InputType: SettingInputType
  Options: string[] | null
}

export interface BasicSetting {
  [key:string]: string | number
}

export interface CategorizedSettings {
  [key:string]: ExtendedSettingDisplay
}

export interface ExtendedSettingDisplay {
  [key:string]: ExtendedSettingDisplayProps & { Value: string | number }
}

export interface ExtendedSettings {
  [key:string]: ExtendedSettingDisplayProps
}