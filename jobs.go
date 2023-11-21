package zoau

import (
	"strings"
	"time"
)

func CancelJob(jobId string, args *CancelJobArgs) error {
	options := make([]string, 0)
	if args != nil {
		if args.Purge {
			options = append(options, "P")
		} else {
			options = append(options, "C")
		}

		if args.JobName != nil {
			options = append(options, *args.JobName)
		} else {
			options = append(options, "*")
		}
	}

	options = append(options, jobId)

	if _, _, err := execZaouCmd("jcan", options); err != nil {
		return err
	}

	duration := time.Second * 10
	if args.Timeout != nil {
		duration = *args.Timeout
	}

	jobChan := make(chan struct {
		Job *Job
		Err error
	}, 1)
	timeOutChan := make(chan struct{}, 1)

	asyncGetJob := func(jobId string, jobChan chan struct {
		Job *Job
		Err error
	},
	) {
		job, err := GetJob(jobId)
		jobChan <- struct {
			Job *Job
			Err error
		}{
			Job: job,
			Err: err,
		}
	}

	go asyncGetJob(jobId, jobChan)
	go timeOut(duration, timeOutChan)

	for {
		select {
		case j := <-jobChan:
			if j.Err != nil || j.Job != nil {
				return j.Err
			}
			go asyncGetJob(jobId, jobChan)
		case <-timeOutChan:
			return nil
		}
	}
}

func GetJob(jobId string) (*Job, error) {
	if out, err := ListingJobs(&jobId, nil); err != nil {
		return nil, err
	} else {
		return &out[0], nil
	}
}

func ListingJobs(jobId *string, jobOwner *string) ([]Job, error) {
	pattern := ""
	if jobId != nil {
		pattern += "/" + *jobId
	} else if jobOwner != nil {
		pattern += "/" + *jobOwner
	}

	stdout, _, err := execZaouCmd("jls", []string{pattern})
	if err != nil {
		return nil, err
	}

	lines := strings.Split(stdout, "\n")

	jobs := make([]Job, len(lines))
	for i, l := range lines {
		output := ParseLine(l)
		jobs[i] = Job{
			Owner:  defaultOrNil(output[0]),
			Name:   defaultOrNil(output[1]),
			Id:     defaultOrNil(output[2]),
			Status: defaultOrNil(output[3]),
			Rc:     defaultOrNil(output[4]),
		}
	}

	return jobs, nil
}

func defaultOrNil(value string) *string {
	if value == "?" {
		return nil
	}
	return &value
}

func ListJobDDs(jobId string, args *JobDDsArgs) ([]JobDDs, error) {
	options := []string{jobId}
	if args != nil {
		pattern := ""
		if args.Owner != nil {
			pattern += "/" + *args.Owner
		}
		if args.Prefix != nil {
			pattern += "/" + *args.Prefix
		}
		options = append(options, pattern)
	}

	stdout, _, err := execZaouCmd("ddls", options)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(stdout, "\n")
	jobsDDs := make([]JobDDs, len(lines))

	for i, l := range lines {
		output := ParseLine(l)
		jobDDs := JobDDs{
			StepName: output[0],
			Dataset:  output[1],
			Format:   output[3],
			Length:   output[4],
			RecNum:   output[5],
		}

		if output[2] != "-" {
			jobDDs.ProcStep = &output[2]
		}

		jobsDDs[i] = jobDDs
	}

	return jobsDDs, nil
}

func ReadJobOutput(jobId string, stepname string, dataset string, args *ReadJobOutputArgs) (string, error) {
	options := []string{jobId, stepname}

	if args != nil {
		if args.ProcStep != nil {
			options = append(options, *args.ProcStep)
		}
	}

	options = append(options, dataset)

	if args != nil {
		pattern := ""
		if args.Owner != nil {
			pattern += "/" + *args.Owner
		}
		if args.Prefix != nil {
			pattern += "/" + *args.Prefix
		}
		options = append(options, pattern)
	}

	return execSimpleStringCmd("pjdd", options)
}

func SubmitJob(dataset string, args *SubmitArgs) (*Job, error) {
	jobId, err := execSimpleStringCmd("jsub", []string{dataset})
	if err != nil {
		return nil, err
	}

	duration := time.Second * 10
	if args != nil {
		if !args.Wait {
			return nil, nil
		}
		if args.Timeout != nil {
			duration = *args.Timeout
		}
	}

	timeOutChan := make(chan struct{}, 0)
	timeOut(duration, timeOutChan)
	<-timeOutChan

	return GetJob(jobId)
}
