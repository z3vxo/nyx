package byte

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Reader struct {
	r   *bytes.Reader
	err error
}

func (r *Reader) Read4() int32 {
	if r.err != nil {
		return 0
	}
	var val int32
	r.err = binary.Read(r.r, binary.LittleEndian, &val)
	return val
}

func (r *Reader) Read1() byte {
	if r.err != nil {
		return 0
	}
	val, err := r.r.ReadByte()
	r.err = err
	return val
}

func (r *Reader) Read2() int16 {
	if r.err != nil {
		return 0
	}
	var val int16
	r.err = binary.Read(r.r, binary.LittleEndian, &val)
	return val
}

func (r *Reader) ReadString(len int32) string {
	if r.err != nil {
		return ""
	}
	buf := make([]byte, len)
	_, r.err = io.ReadFull(r.r, buf)
	return string(buf)

}

/*   Agent Register wire format
 * [GUID LEN] 4 BYTES
 * [GUID STR] N BYTES
 * ->
 * [USERNAME LEN] 4 BYTES
 * [USERNAME STR] N BYTES
 * ->
 * [HOSTNAME LEN] 4 BYTES
 * [HOSTNAME STR] N BYTES
 * ->
 * [INTERNAL IP LEN] 4 BYTES
 * [INTERNAL IP STR] N BYTES
 * ->
 * [PROCESS_PATH LEN] 4 BYTES
 * [PROCESS_PATH STR] N BYTES
 * ->
 * [PID] 4 BYTES
 * -->
 * [PPID] 4 BYTES
 * ->
 * [IS ELEVATED] 1 BYTE
 * ->
 * [ARCH] 1 BYTE | 1 == x64, 0 == x86
 * ->
 * [MINOR VERSION] 2 BYTES
 * [MAJOR VERSION] 2 BYTES
 * [BUILD VERSION] 2 BYTES
 *
 */

type ClientRegister struct {
	Guid       string
	User       string
	Host       string
	InternaIP  string
	ExternalIP string
	ProcPath   string
	Pid        int32
	Ppid       int32
	IsElev     byte
	Arch       byte
	Minor      int16
	Major      int16
	Build      int16
}

func ExtractRegistrationDetails(IP string, r *bytes.Reader) (ClientRegister, error) {
	rd := &Reader{r: r}

	guid := rd.ReadString(rd.Read4())
	Username := rd.ReadString(rd.Read4())
	Hostname := rd.ReadString(rd.Read4())
	InternalIP := rd.ReadString(rd.Read4())
	ProcessPath := rd.ReadString(rd.Read4())
	Pid := rd.Read4()
	PPid := rd.Read4()
	IsElev := rd.Read1()
	Arch := rd.Read1()
	Minor := rd.Read2()
	Major := rd.Read2()
	BuildVer := rd.Read2()
	if rd.err != nil {
		return ClientRegister{}, rd.err
	}

	Res := ClientRegister{
		Guid:       guid,
		User:       Username,
		Host:       Hostname,
		InternaIP:  InternalIP,
		ExternalIP: IP,
		ProcPath:   ProcessPath,
		Pid:        Pid,
		Ppid:       PPid,
		IsElev:     IsElev,
		Arch:       Arch,
		Minor:      Minor,
		Major:      Major,
		Build:      BuildVer,
	}

	return Res, nil
}
