package resources

import "fmt"

func (t Resource) aws_cloudwatch_metric_alarm_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe cloudwatch_alarm(%v) do\n"+
		"  it { should exist }\n" +
		"  its(:metric_name) { should eq %v }\n" +
		"  its(:alarm_description) { should eq %v }\n" +
		"  its(:namespace) { should eq %v }\n" +
		"  its(:actions_enabled) { should eq %v }\n" +
		"  its(:comparison_operator) { should eq %v }\n" +
		"  its(:threshold) { should eq %v }\n" +
		"  its(:evaluation_periods) { should eq %v }\n" +
		"  its(:unit) { should eq %v }\n" +
		"  its(:period) { should eq %v }\n" +
		"  its(:statistic) { should eq %v }\n" +
		"  its(:extended_statistic) { should eq %v }\n" +
		"  its(:evaluate_low_sample_count_percentile) { should eq %v }\n" +
		"  its(:treat_missing_data) { should eq %v }\n" +
		t.alarm_attrs() +
		"end\n",
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"alarm_name").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"metric_name").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"alarm_description").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"namespace").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"actions_enabled").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"comparison_operator").String()),
		get_attribute_by_name(t.Attrs,"threshold").Uint(),
		get_attribute_by_name(t.Attrs,"evaluation_periods").Uint(),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"unit").String()),
		get_attribute_by_name(t.Attrs,"period").Uint(),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"statistic").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"extended_statistic").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"evaluate_low_sample_count_percentiles").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"treat_missing_data").String()),
	)
}

func (t Resource) alarm_attrs() (ret string) {
	ok_actions := get_list_items_by_pattern(t.Attrs, `^ok_actions\.[0-9]+$`)
	alarm_actions := get_list_items_by_pattern(t.Attrs, `^alarm_actions\.[0-9]+$`)
	insufficient_data_actions := get_list_items_by_pattern(t.Attrs, `^insufficient_data_actions\.[0-9]+$`)
	ret += create_ruby_string_array(ok_actions, "  its(:ok_actions) { should =~ [%v]}\n")
	ret += create_ruby_string_array(alarm_actions, "  its(:alarm_actions) { should =~ [%v]}\n")
	ret += create_ruby_string_array(insufficient_data_actions, "  its(:insufficient_data_actions) { should =~ [%v]}\n")
	return
}
