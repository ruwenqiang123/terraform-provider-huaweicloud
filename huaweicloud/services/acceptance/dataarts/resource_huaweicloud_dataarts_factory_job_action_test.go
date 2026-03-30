package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFactoryJobAction_basic(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceName()

		actionWithSequence   = "huaweicloud_dataarts_factory_job_action.action_with_sequence"
		rcActionWithSequence = acceptance.InitResourceCheck(actionWithSequence, &obj, getFactoryJobResourceFunc)

		startWithoutParams   = "huaweicloud_dataarts_factory_job_action.start_without_params"
		rcStartWithoutParams = acceptance.InitResourceCheck(startWithoutParams, &obj, getFactoryJobResourceFunc)

		stop   = "huaweicloud_dataarts_factory_job_action.stop"
		rcStop = acceptance.InitResourceCheck(stop, &obj, getFactoryJobResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcActionWithSequence.CheckResourceDestroy(),
			rcStartWithoutParams.CheckResourceDestroy(),
			rcStop.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testFactoryJobAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					// Check the first job is started successfully.
					rcActionWithSequence.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(actionWithSequence, "job_name", "huaweicloud_dataarts_factory_job.test.0", "name"),
					resource.TestCheckResourceAttr(actionWithSequence, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(actionWithSequence, "action", "start"),
					resource.TestCheckResourceAttr(actionWithSequence, "status", "NORMAL"),
					// Check the second job is started successfully.
					rcStartWithoutParams.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(startWithoutParams, "job_name", "huaweicloud_dataarts_factory_job.test.1", "name"),
					resource.TestCheckResourceAttr(startWithoutParams, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(startWithoutParams, "action", "start"),
					resource.TestCheckResourceAttr(startWithoutParams, "status", "NORMAL"),
				),
			},
			{
				Config: testFactoryJobAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					// Check the first job is started successfully.
					rcActionWithSequence.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(actionWithSequence, "job_name", "huaweicloud_dataarts_factory_job.test.0", "name"),
					resource.TestCheckResourceAttr(actionWithSequence, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(actionWithSequence, "action", "stop"),
					resource.TestCheckResourceAttr(actionWithSequence, "status", "STOPPED"),
					// Check if the one-time resource used to start the second job exists.
					rcStartWithoutParams.CheckResourceExists(),
					// Check the second job is stopped successfully.
					rcStop.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(stop, "job_name", "huaweicloud_dataarts_factory_job.test.1", "name"),
					resource.TestCheckResourceAttr(stop, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(stop, "action", "stop"),
					resource.TestCheckResourceAttr(stop, "status", "STOPPED"),
				),
			},
		},
	})
}

func testFactoryJobAction_realTime_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dataarts_factory_job" "test" {
  count = 2

  name         = "%[1]s_${count.index}"
  workspace_id = "%[2]s"
  process_type = "REAL_TIME"

  nodes {
    name = "SMN_%[1]s_${count.index}"
    type = "SMN"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "topic"
      value = huaweicloud_smn_topic.test.topic_urn
    }

    properties {
      name  = "messageType"
      value = "NORMAL"
    }

    properties {
      name  = "message"
      value = "terraform acceptance test"
    }
  }

  schedule {
    type = "EXECUTE_ONCE"
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testFactoryJobAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# The first job is started using the start_date, ignore_first_self_dep, and job_params parameters.
# Subsequent steps will stop this (the first) job.
resource "huaweicloud_dataarts_factory_job_action" "action_with_sequence" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test[0].name
  process_type = huaweicloud_dataarts_factory_job.test[0].process_type

  # Obtain the start date 48 hours from now, the format is YYmmDD, such as '20060102'
  start_date            = formatdate("YYmmDD", timeadd(timestamp(), "48h"))
  ignore_first_self_dep = true

  dynamic "job_params" {
    for_each = try(huaweicloud_dataarts_factory_job.test[0].nodes[0].properties, [])

    content {
      name  = job_params.value.name
      value = job_params.value.value
    }
  }

  lifecycle {
    ignore_changes = [
      start_date,
    ]
  }
}

# Create a new one-time action resource to start the second job.
resource "huaweicloud_dataarts_factory_job_action" "start_without_params" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test[1].name
  process_type = huaweicloud_dataarts_factory_job.test[1].process_type
}
`, testFactoryJobAction_realTime_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testFactoryJobAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

# Following the previous step, stop the first job that has already started by changing the Terraform resources.
resource "huaweicloud_dataarts_factory_job_action" "action_with_sequence" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "stop"
  job_name     = huaweicloud_dataarts_factory_job.test[0].name
  process_type = huaweicloud_dataarts_factory_job.test[0].process_type
}

# Retain this one-time action resource.
resource "huaweicloud_dataarts_factory_job_action" "start_without_params" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test[1].name
  process_type = huaweicloud_dataarts_factory_job.test[1].process_type
}

# Create a new one-time action resource to stop the second job.
resource "huaweicloud_dataarts_factory_job_action" "stop" {
  depends_on = [
    huaweicloud_dataarts_factory_job_action.start_without_params,
  ]

  workspace_id = "%[2]s"
  action       = "stop"
  job_name     = huaweicloud_dataarts_factory_job.test[1].name
  process_type = huaweicloud_dataarts_factory_job.test[1].process_type
}
`, testFactoryJobAction_realTime_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func TestAccFactoryJobAction_batchJob(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceName()

		actionWithSequence   = "huaweicloud_dataarts_factory_job_action.action_with_sequence"
		rcActionWithSequence = acceptance.InitResourceCheck(actionWithSequence, &obj, getFactoryJobResourceFunc)

		startWithoutParams   = "huaweicloud_dataarts_factory_job_action.start_without_params"
		rcStartWithoutParams = acceptance.InitResourceCheck(startWithoutParams, &obj, getFactoryJobResourceFunc)

		startImmediately   = "huaweicloud_dataarts_factory_job_action.start_immediately"
		rcStartImmediately = acceptance.InitResourceCheck(startImmediately, &obj, getFactoryJobResourceFunc)

		stopStart   = "huaweicloud_dataarts_factory_job_action.stop"
		rcStopStart = acceptance.InitResourceCheck(stopStart, &obj, getFactoryJobResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsCdmName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcActionWithSequence.CheckResourceDestroy(),
			rcStartWithoutParams.CheckResourceDestroy(),
			rcStartImmediately.CheckResourceDestroy(),
			rcStopStart.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testFactoryJobAction_batchPinelineJob_step1(name),
				Check: resource.ComposeTestCheckFunc(
					// Check the first job is started successfully.
					rcActionWithSequence.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(actionWithSequence, "job_name", "huaweicloud_dataarts_factory_job.test.0", "name"),
					resource.TestCheckResourceAttr(actionWithSequence, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(actionWithSequence, "action", "start"),
					resource.TestCheckResourceAttr(actionWithSequence, "status", "SCHEDULING"),
					// Check the second job is started successfully.
					rcStartWithoutParams.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(startWithoutParams, "job_name", "huaweicloud_dataarts_factory_job.test.1", "name"),
					resource.TestCheckResourceAttr(startWithoutParams, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(startWithoutParams, "action", "start"),
					resource.TestCheckResourceAttr(startWithoutParams, "status", "SCHEDULING"),
					// Check the the test of third job is started successfully (immediately).
					rcStartImmediately.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(startImmediately, "job_name", "huaweicloud_dataarts_factory_job.test.2", "name"),
					resource.TestCheckResourceAttr(startImmediately, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(startImmediately, "action", "run-immediate"),
					resource.TestMatchResourceAttr(startImmediately, "instance_status", regexp.MustCompile("^(success|fail)$")),
				),
			},
			{
				Config: testFactoryJobAction_batchPinelineJob_step2(name),
				Check: resource.ComposeTestCheckFunc(
					// Check the first job is started successfully.
					rcActionWithSequence.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(actionWithSequence, "job_name", "huaweicloud_dataarts_factory_job.test.0", "name"),
					resource.TestCheckResourceAttr(actionWithSequence, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(actionWithSequence, "action", "stop"),
					resource.TestCheckResourceAttr(actionWithSequence, "status", "STOPPED"),
					// Check if the one-time resource used to start the second job exists.
					rcStartWithoutParams.CheckResourceExists(),
					// Check the second job is stopped successfully.
					rcStopStart.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(stopStart, "job_name", "huaweicloud_dataarts_factory_job.test.1", "name"),
					resource.TestCheckResourceAttr(stopStart, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(stopStart, "action", "stop"),
					resource.TestCheckResourceAttr(stopStart, "status", "STOPPED"),
				),
			},
		},
	})
}

func testFactoryJobAction_batchPinelineJob_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_job" "test" {
  count = 3

  name         = "%[1]s_${count.index}"
  workspace_id = "%[2]s"
  process_type = "BATCH"

  nodes {
    name = "Rest_client_%[1]s_${count.index}"
    type = "RESTAPI"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "url"
      value = "https://www.huaweicloud.com/"
    }

    properties {
      name  = "method"
      value = "GET"
    }

    properties {
      name  = "retry"
      value = "false"
    }

    properties {
      name  = "requestMode"
      value = "sync"
    }

    properties {
      name  = "securityAuthentication"
      value = "NONE"
    }

    properties {
      name  = "agentName"
      value = "%[3]s"
    }
  }

  schedule {
    type = "CRON"
    cron {
      expression = "0 0 0 * * ?"
      start_time = "2024-07-24T16:14:04+08"
    }
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CDM_NAME)
}

func testFactoryJobAction_batchPinelineJob_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# The first job is started using the start_date, ignore_first_self_dep, and job_params parameters.
# Subsequent steps will stop this (the first) job.
resource "huaweicloud_dataarts_factory_job_action" "action_with_sequence" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test[0].name
  process_type = huaweicloud_dataarts_factory_job.test[0].process_type

  # Obtain the start date 48 hours from now, the format is YYmmDD, such as '20060102'
  start_date            = formatdate("YYmmDD", timeadd(timestamp(), "48h"))
  ignore_first_self_dep = true

  dynamic "job_params" {
    for_each = try(huaweicloud_dataarts_factory_job.test[0].nodes[0].properties, [])

    content {
      name  = job_params.value.name
      value = job_params.value.value
    }
  }

  lifecycle {
    ignore_changes = [
      start_date,
    ]
  }
}

# Create a new one-time action resource to start the second job.
resource "huaweicloud_dataarts_factory_job_action" "start_without_params" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test[1].name
  process_type = huaweicloud_dataarts_factory_job.test[1].process_type
}

# Create a new one-time action resource to start the third job immediately.
resource "huaweicloud_dataarts_factory_job_action" "start_immediately" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id       = "%[2]s"
  action             = "run-immediate"
  job_name           = huaweicloud_dataarts_factory_job.test[2].name
  process_type       = huaweicloud_dataarts_factory_job.test[2].process_type
  use_execution_user = "true"

  dynamic "job_params" {
    for_each = try(huaweicloud_dataarts_factory_job.test[0].nodes[0].properties, [])

    content {
      name  = job_params.value.name
      value = job_params.value.value
    }
  }
}
`, testFactoryJobAction_batchPinelineJob_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testFactoryJobAction_batchPinelineJob_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

# Following the previous step, stop the first job that has already started by changing the Terraform resources.
resource "huaweicloud_dataarts_factory_job_action" "action_with_sequence" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "stop"
  job_name     = huaweicloud_dataarts_factory_job.test[0].name
  process_type = huaweicloud_dataarts_factory_job.test[0].process_type
}

# Retain this one-time action resource.
resource "huaweicloud_dataarts_factory_job_action" "start_without_params" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test[1].name
  process_type = huaweicloud_dataarts_factory_job.test[1].process_type
}

# Create a new one-time action resource to stop the second job.
resource "huaweicloud_dataarts_factory_job_action" "stop" {
  depends_on = [
    huaweicloud_dataarts_factory_job_action.start_without_params,
  ]

  workspace_id = "%[2]s"
  action       = "stop"
  job_name     = huaweicloud_dataarts_factory_job.test[1].name
  process_type = huaweicloud_dataarts_factory_job.test[1].process_type
}

# Retain this one-time action resource.
resource "huaweicloud_dataarts_factory_job_action" "start_immediately" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test,
  ]

  workspace_id       = "%[2]s"
  action             = "run-immediate"
  job_name           = huaweicloud_dataarts_factory_job.test[2].name
  process_type       = huaweicloud_dataarts_factory_job.test[2].process_type
  use_execution_user = "true"

  dynamic "job_params" {
    for_each = try(huaweicloud_dataarts_factory_job.test[0].nodes[0].properties, [])

    content {
      name  = job_params.value.name
      value = job_params.value.value
    }
  }
}
`, testFactoryJobAction_batchPinelineJob_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
